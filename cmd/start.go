package cmd

import (
	"github.com/omeid/livereload"
	"github.com/oneut/lrp/command"
	"github.com/oneut/lrp/config"
	"github.com/oneut/lrp/monitor"
	"github.com/oneut/lrp/proxy"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start live reload proxy",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func run() {
	log.Info("Start live reload proxy")
	lrp := &LivereloadProxy{
		Config: config.GetConfig(),
	}

	lrp.startTasks()
	lrp.startLivereload()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	for {
		select {
		case <-sigChan:
			// @todo terminate処理
			os.Exit(0)
		}
	}
}

type LivereloadProxy struct {
	Config     *config.Config
	Tasks      map[string]*Task
	Livereload *livereload.Server
}

type Task struct {
	Command *command.Command
	Monitor monitor.Monitorer
}

func (lrp *LivereloadProxy) startTasks() {
	for name, task := range lrp.Config.Tasks {
		go func() {
			lrp.startTask(name, task)
		}()
	}
}

func (lrp *LivereloadProxy) startTask(name string, taskConfig config.Task) {

	c := command.NewCommand(name, taskConfig.Command)
	m := monitor.NewMonitor(name, taskConfig.Monitor)

	fn := func(message string) {
		c.Restart()
		lrp.Livereload.Reload(message, true)
	}

	c.Run(fn)
	m.Run(fn)

	lrp.Tasks[name] = &Task{
		Command: c,
		Monitor: m,
	}
}

func (lrp *LivereloadProxy) startLivereload() {
	lrp.Livereload = livereload.New("LivereloadProxy")

	r := &RegexpHandler{}

	scriptPath := "/livereload.js"
	r.Handle("/livereload", lrp.Livereload)
	r.HandleFunc(scriptPath, livereload.LivereloadScript)
	r.HandleFunc("*", lrp.handler)

	http.Handle("/", r)
	go func() {
		defer lrp.Livereload.Close()
		http.ListenAndServe(lrp.Config.ProxyHost, nil)
	}()
}

func (lrp *LivereloadProxy) handler(w http.ResponseWriter, r *http.Request) {
	scriptPath := "/livereload.js"
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = lrp.Config.SourceHost
	}

	modifier := func(res *http.Response) error {
		//contentType := res.Header.Get("Content-type")
		//log.Info(contentType)
		//if contentType != "text/html" {
		//	return nil
		//}
		proxyDocument := &proxy.ProxyDocument{res.Body}
		buf := proxyDocument.CreateBytesBufferWithLiveReloadScriptPath(scriptPath)
		s := buf.String()
		s = strings.Replace(s, lrp.Config.SourceHost, lrp.Config.ProxyHost, -1)
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

type RetryTransport struct {
}

func (rt *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for {
		res, err := http.DefaultTransport.RoundTrip(req)
		if err == nil {
			return res, nil
		}

		time.Sleep(200 * time.Millisecond)
	}
}

type route struct {
	pattern string
	handler http.Handler
}

type RegexpHandler struct {
	routes []*route
}

func (h *RegexpHandler) Handle(pattern string, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, handler})
}

func (h *RegexpHandler) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern == "*" {
			route.handler.ServeHTTP(w, r)
			return
		}
		if route.pattern == r.URL.Path {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}
