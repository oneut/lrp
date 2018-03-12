package livereloadproxy

import (
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	"github.com/omeid/livereload"
)

func NewProxy(proxyHost string, sourceHost string) *Proxy {
	return &Proxy{
		livereload: livereload.New("LivereloadProxy"),
		proxyHost:  proxyHost,
		sourceHost: sourceHost,
	}
}

type Proxy struct {
	livereload *livereload.Server
	proxyHost  string
	sourceHost string
}

func (p *Proxy) Run() {
	r := NewRouter()
	scriptPath := "/livereload.js"
	r.Handle("/livereload", p.livereload)
	r.HandleFunc(scriptPath, livereload.LivereloadScript)
	r.HandleFunc("*", func(w http.ResponseWriter, r *http.Request) {
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
			buf := proxyBody.getBytesBufferWithLiveReloadScriptPath(scriptPath)
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
	})
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
