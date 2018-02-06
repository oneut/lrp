package main

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	"github.com/omeid/livereload"

	"github.com/oneut/llp/monitor"
	"github.com/oneut/llp/proxy"
)

func main() {
	log.Info("Start Livereload Proxy")
	config := getConfig()
	llp := &LivereloadProxy{}
	for name, task := range config.Tasks {
		go func() {
			log.WithFields(log.Fields{
				"name": name,
			}).Info("Start monitor")
			llp.StartMonitor(name, task)
		}()
	}
	llp.StartLivereload()
}

type LivereloadProxy struct {
	Monitor    map[string]*monitor.Monitor
	Livereload *livereload.Server
}

func (llp *LivereloadProxy) StartMonitor(name string, task monitor.Task) {
	m := monitor.NewMonitor(name, task)
	m.Run(func(message string) {
		log.Info("Reload")
		llp.Livereload.Reload(message, true)
	})
	llp.Monitor[name] = m
}

func (llp *LivereloadProxy) StartLivereload() {
	llp.Livereload = livereload.New("LivereloadProxy")

	scriptPath := "/livereload.js"
	http.Handle("/livereload", llp.Livereload)
	http.HandleFunc(scriptPath, livereload.LivereloadScript)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		director := func(req *http.Request) {
			// @todo
			req.URL.Scheme = "http"
			req.URL.Host = ":8080"
		}

		modifier := func(res *http.Response) error {
			proxyDocument := &proxy.ProxyDocument{res.Body}
			buf := proxyDocument.CreateBytesBufferWithLiveReloadScriptPath(scriptPath)

			res.Header.Set("Content-Length", strconv.Itoa(buf.Len()))
			res.Body = ioutil.NopCloser(strings.NewReader(buf.String()))
			return nil
		}

		proxy := &httputil.ReverseProxy{
			Director:       director,
			ModifyResponse: modifier,
		}
		proxy.ServeHTTP(w, r)
	})

	// @todo
	http.ListenAndServe(":9000", nil)
}
