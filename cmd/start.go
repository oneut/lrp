package cmd

import (
	"github.com/omeid/livereload"
	"github.com/oneut/lrp/config"
	"github.com/oneut/lrp/monitor"
	"github.com/oneut/lrp/proxy"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start live reload proxy",
	Long:  `A longer description that spans multiple lines and likely contains examples and usage of using your command. For example: Cobra is a CLI library for Go that empowers applications.This application is a tool to generate the needed files to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func run() {
	log.Info("Start live reload proxy")
	cfg := config.GetConfig()
	lrp := &LivereloadProxy{}
	for name, task := range cfg.Tasks {
		go func() {
			log.WithFields(log.Fields{
				"name": name,
			}).Info("Start monitor")
			lrp.StartMonitor(name, task)
		}()
	}
	lrp.StartLivereload(cfg.ProxyHost, cfg.SourceHost)
}

type LivereloadProxy struct {
	Monitor    map[string]*monitor.Monitor
	Livereload *livereload.Server
}

func (llp *LivereloadProxy) StartMonitor(name string, task monitor.Task) {
	m := monitor.NewMonitor(name, task)
	m.Run(func(message string) {
		llp.Livereload.Reload(message, true)
	})
	llp.Monitor[name] = m
}

func (llp *LivereloadProxy) StartLivereload(proxyHost string, sourceHost string) {
	llp.Livereload = livereload.New("LivereloadProxy")

	scriptPath := "/livereload.js"
	http.Handle("/livereload", llp.Livereload)
	http.HandleFunc(scriptPath, livereload.LivereloadScript)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		director := func(req *http.Request) {
			// @todo
			req.URL.Scheme = "http"
			req.URL.Host = sourceHost
		}

		modifier := func(res *http.Response) error {
			proxyDocument := &proxy.ProxyDocument{res.Body}
			buf := proxyDocument.CreateBytesBufferWithLiveReloadScriptPath(scriptPath)
			s := buf.String()
			s = strings.Replace(s, sourceHost, proxyHost, -1)
			res.Header.Set("Content-Length", strconv.Itoa(len(s)))
			res.Body = ioutil.NopCloser(strings.NewReader(s))
			return nil
		}

		proxy := &httputil.ReverseProxy{
			Director:       director,
			ModifyResponse: modifier,
			Transport:      &RetryTransport{},
		}
		proxy.ServeHTTP(w, r)
	})

	// @todo
	http.ListenAndServe(proxyHost, nil)
}

type RetryTransport struct {
}

func (rt *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for {
		res, err := http.DefaultTransport.RoundTrip(req)
		if err == nil {
			log.Printf("%d\t%s\t%s\n", res.StatusCode, req.Method, req.URL.String())
			return res, err
		}

		time.Sleep(200 * time.Millisecond)
	}
}
