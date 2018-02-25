package cmd

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/omeid/livereload"
	"github.com/oneut/lrp/command"
	"github.com/oneut/lrp/config"
	"github.com/oneut/lrp/livereloadproxy"
	"github.com/oneut/lrp/monitor"
	"github.com/oneut/lrp/proxy"
	"github.com/spf13/cobra"
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
	log.Println("Start live reload proxy")
	lrp := &LivereloadProxy{
		Config: config.GetConfig(),
		Tasks:  make(map[string]*Task),
	}

	lrp.startTasks()
	lrp.startLivereload()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	for {
		select {
		case <-sigChan:
			lrp.stopTasks()
			lrp.stopLivereload()
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
	Commands map[string]*command.Command
	Monitor  monitor.Monitorer
}

func (lrp *LivereloadProxy) startTasks() {
	for name, task := range lrp.Config.Tasks {
		lrp.startTask(name, task)
	}
}

func (lrp *LivereloadProxy) startTask(name string, taskConfig config.Task) {
	m := monitor.NewMonitor(name, taskConfig.Monitor)

	cmds := make(map[string]*command.Command)
	for cmdName, cmdConfig := range taskConfig.Commands {
		cmds[cmdName] = command.NewCommand(name, cmdName, cmdConfig)
	}

	isReloading := false
	fn := func(message string) {
		if isReloading {
			return
		}

		isReloading = true
		go func() {
			time.Sleep(taskConfig.GetAggregateTimeout())
			for _, cmd := range cmds {
				cmd.Restart()
			}
			lrp.Livereload.Reload(message, true)
			isReloading = false
		}()
	}

	lrp.Tasks[name] = &Task{
		Commands: cmds,
		Monitor:  m,
	}

	for _, cmd := range cmds {
		go cmd.Run(fn)
	}
	go m.Run(fn)
}

func (lrp *LivereloadProxy) startLivereload() {
	lrp.Livereload = livereload.New("LivereloadProxy")

	r := livereloadproxy.NewRouter()
	scriptPath := "/livereload.js"
	r.Handle("/livereload", lrp.Livereload)
	r.HandleFunc(scriptPath, livereload.LivereloadScript)
	r.HandleFunc("*", func(w http.ResponseWriter, r *http.Request) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = lrp.Config.SourceHost
		}

		modifier := func(res *http.Response) error {
			contentType := res.Header.Get("Content-type")
			if !(strings.Contains(contentType, "text/html")) {
				return nil
			}
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
	})
	go func() {
		defer lrp.Livereload.Close()
		http.ListenAndServe(lrp.Config.ProxyHost, r)
	}()
}

func (lrp *LivereloadProxy) stopTasks() {
	for name, _ := range lrp.Config.Tasks {
		lrp.stopTask(name)
	}
}

func (lrp *LivereloadProxy) stopTask(name string) {
	task := lrp.Tasks[name]
	for _, cmd := range task.Commands {
		cmd.Stop()
	}
	task.Monitor.Stop()
}

func (lrp *LivereloadProxy) stopLivereload() {
	lrp.Livereload.Close()
}

type RetryTransport struct {
}

func (rt *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for {
		res, err := http.DefaultTransport.RoundTrip(req)
		if err == nil {
			return res, nil
		}

		// Retry
		time.Sleep(500 * time.Millisecond)
	}
}
