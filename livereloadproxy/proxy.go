package livereloadproxy

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"path"
	"strconv"
	"strings"

	"github.com/omeid/livereload"
)

func NewProxy(proxyHost string, staticPath string, sourceHost string) *Proxy {
	return &Proxy{
		livereload: livereload.New("LivereloadProxy"),
		proxyHost:  proxyHost,
		sourceHost: sourceHost,
		staticPath: staticPath,
		scriptPath: "/livereload.js",
	}
}

type Proxy struct {
	livereload *livereload.Server
	proxyHost  string
	sourceHost string
	scriptPath string
	staticPath string
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
		http.ListenAndServe(p.proxyHost, r)
	}()
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
		req.URL.Scheme = "http"
		req.URL.Host = p.sourceHost
	}

	modifier := func(res *http.Response) error {
		contentType := res.Header.Get("Content-type")
		if !(strings.Contains(contentType, "text/html")) {
			return nil
		}
		proxyBody := &ProxyBody{res.Body}
		buf := proxyBody.getBytesBufferWithLiveReloadScriptPath(p.scriptPath)
		s := buf.String()
		s = strings.Replace(s, p.sourceHost, p.proxyHost, -1)
		res.Header.Set("Content-Length", strconv.Itoa(len(s)))
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
	if p.sourceHost == "" {
		return false
	}

	return true
}
