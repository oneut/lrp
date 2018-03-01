package livereloadproxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	"github.com/omeid/livereload"
	"github.com/oneut/lrp/command"
	"github.com/oneut/lrp/config"
	"github.com/oneut/lrp/monitor"
)

func NewLivereloadProxy() *LivereloadProxy {
	return &LivereloadProxy{
		Config: config.GetConfig(),
		Tasks:  make(map[string]*Task),
	}
}

type LivereloadProxy struct {
	Config     *config.Config
	Tasks      map[string]*Task
	Livereload *livereload.Server
}

func (lrp *LivereloadProxy) Run() {
	lrp.startTasks()
	lrp.startLivereload()
}

func (lrp *LivereloadProxy) startTasks() {
	for name, task := range lrp.Config.Tasks {
		lrp.startTask(name, task)
	}
}

func (lrp *LivereloadProxy) startTask(name string, taskConfig config.Task) {
	m := monitor.NewMonitor(name, taskConfig.Monitor)

	cmds := make(map[string]command.Commander)
	for cmdName, cmdConfig := range taskConfig.Commands {
		cmds[cmdName] = command.NewCommand(name, cmdName, cmdConfig)
		fmt.Printf("%#v", cmds[cmdName])
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

	r := NewRouter()
	scriptPath := "/livereload.js"
	r.Handle("/livereload", lrp.Livereload)
	r.HandleFunc(scriptPath, livereload.LivereloadScript)
	r.HandleFunc("*", func(w http.ResponseWriter, r *http.Request) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = lrp.Config.GetSourceHost()
		}

		modifier := func(res *http.Response) error {
			contentType := res.Header.Get("Content-type")
			if !(strings.Contains(contentType, "text/html")) {
				return nil
			}
			proxyBody := &ProxyBody{res.Body}
			buf := proxyBody.CreateBytesBufferWithLiveReloadScriptPath(scriptPath)
			s := buf.String()
			s = strings.Replace(s, lrp.Config.GetSourceHost(), lrp.Config.GetProxyHost(), -1)
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
		http.ListenAndServe(lrp.Config.GetProxyHost(), r)
	}()
}

func (lrp *LivereloadProxy) Stop() {
	lrp.stopTasks()
	lrp.stopLivereload()
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
