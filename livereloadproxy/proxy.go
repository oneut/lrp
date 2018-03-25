package livereloadproxy

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"

	"github.com/omeid/livereload"
	"github.com/oneut/lrp/config"
	"github.com/skratchdot/open-golang/open"
)

func NewProxy(proxyConfig config.Proxy, sourceConfig config.Source) *Proxy {
	proxy := &Proxy{
		livereload: livereload.New("LivereloadProxy"),
		proxyURL: &url.URL{
			Scheme: proxyConfig.GetScheme(),
			Host:   proxyConfig.Host,
		},
		staticPath:    proxyConfig.StaticPath,
		isBrowserOpen: proxyConfig.IsBrowserOpen(),
		sourceURL: &url.URL{
			Scheme: sourceConfig.GetScheme(),
			Host:   sourceConfig.Host,
		},
		scriptPath: "/livereload.js",
	}

	for _, replace := range sourceConfig.Replaces {
		if !(replace.IsValid()) {
			continue
		}

		proxy.AddSourceReplace(replace)
	}

	return proxy
}

type Proxy struct {
	livereload      *livereload.Server
	proxyURL        *url.URL
	scriptPath      string
	isBrowserOpen   bool
	sourceURL       *url.URL
	staticPath      string
	sourceReplacers []SourceReplacer
}

func (p *Proxy) AddSourceReplace(replaceConfig config.Replace) {
	p.sourceReplacers = append(p.sourceReplacers, NewSourceReplacer(replaceConfig))
}

func (p *Proxy) Run() {
	r := NewRouter()
	r.Handle("/livereload", p.livereload)
	r.HandleFunc(p.scriptPath, livereload.LivereloadScript)

	if p.hasStaticPath() {
		fs := http.Dir(p.staticPath)
		r.HandleFunc("*", func(w http.ResponseWriter, r *http.Request) {
			p.handleStatic(fs, w, r)
		})
	} else if p.hasSourceHost() {
		r.HandleFunc("*", p.handleReverseProxy)
	}

	go func() {
		defer p.livereload.Close()
		http.ListenAndServe(p.proxyURL.Host, r)
	}()

	if p.isBrowserOpen {
		err := open.Start(p.proxyURL.String())
		if err != nil {
			panic(err)
		}
	}
}

func (p *Proxy) Close() {
	p.livereload.Close()
}

func (p *Proxy) Reload(message string) {
	p.livereload.Reload(message, true)
}

func (p *Proxy) handleStatic(fs http.FileSystem, w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}

	name := path.Clean(upath)

	const indexPage = "/index.html"

	if strings.HasSuffix(r.URL.Path, indexPage) {
		http.Redirect(w, r, "./", http.StatusFound)
		return
	}

	f, err := fs.Open(name)
	if err != nil {
		if !(p.hasSourceHost()) {
			panic(err)
		}

		p.handleReverseProxy(w, r)
		return
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		panic(err)
	}

	// use contents of index.html for directory, if present
	if d.IsDir() {
		index := strings.TrimSuffix(name, "/") + indexPage
		ff, err := fs.Open(index)

		if err != nil {
			if !(p.hasSourceHost()) {
				panic(err)
			}

			p.handleReverseProxy(w, r)
			return
		}

		defer ff.Close()
		dd, err := ff.Stat()
		if err != nil {
			if !(p.hasSourceHost()) {
				panic(err)
			}

			p.handleReverseProxy(w, r)
			return
		}

		name = index
		d = dd
		f = ff
	}

	buf := &bytes.Buffer{}
	io.Copy(buf, f)
	contentType := http.DetectContentType(buf.Bytes())
	var reader *bytes.Reader
	if strings.Contains(contentType, "text/html") {
		proxyBody := &ProxyBody{
			ioutil.NopCloser(bytes.NewReader(buf.Bytes())),
		}
		convertedBuf := proxyBody.getBytesBufferWithLiveReloadScriptPath(p.scriptPath)
		reader = bytes.NewReader(convertedBuf.Bytes())
	} else {
		reader = bytes.NewReader(buf.Bytes())
	}

	http.ServeContent(w, r, d.Name(), d.ModTime(), reader)
}

func (p *Proxy) handleReverseProxy(w http.ResponseWriter, r *http.Request) {
	director := func(req *http.Request) {
		req.URL.Scheme = p.sourceURL.Scheme
		req.URL.Host = p.sourceURL.Host
		req.Host = p.sourceURL.Host
		req.Header.Del("Accept-Encoding")
	}

	modifier := func(res *http.Response) error {
		res.Header.Del("Content-Length")
		res.Header.Del("Content-Encoding")
		res.Header.Del("Content-Security-Policy")
		res.Header.Set("Cache-Control", "no-store")

		contentType := res.Header.Get("Content-type")
		if !(strings.Contains(contentType, "text/html")) {
			return nil
		}

		defer res.Body.Close()
		proxyBody := &ProxyBody{res.Body}
		buf := proxyBody.getBytesBufferWithLiveReloadScriptPath(p.scriptPath)

		s := buf.String()

		proxySchemeHost := p.proxyURL.String()
		sourceSchemeHost := p.sourceURL.String()
		sourceHost := (&url.URL{
			Host: p.sourceURL.Host,
		}).String()

		s = strings.Replace(s, sourceSchemeHost, proxySchemeHost, -1)
		s = strings.Replace(s, sourceHost, proxySchemeHost, -1)
		for _, sourceReplacer := range p.sourceReplacers {
			s = sourceReplacer.Replace(s)
		}

		res.Body = ioutil.NopCloser(strings.NewReader(s))
		return nil
	}

	rp := &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifier,
		Transport:      &RetryTransport{},
	}
	rp.ServeHTTP(w, r)
}

func (p *Proxy) hasStaticPath() bool {
	if p.staticPath == "" {
		return false
	}

	return true
}

func (p *Proxy) hasSourceHost() bool {
	if p.sourceURL.Host == "" {
		return false
	}

	return true
}
