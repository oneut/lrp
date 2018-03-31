package livereloadproxy

import (
	"github.com/oneut/lrp/config"
)

func NewLivereloadProxy(cfg config.Config) *LivereloadProxy {
	lrp := &LivereloadProxy{}

	lrp.SetProxy(
		cfg.Proxy,
		cfg.Source,
	)

	for name, taskConfig := range cfg.Tasks {
		lrp.AddTask(name, taskConfig)
	}

	return lrp
}

type LivereloadProxy struct {
	tasks []*Task
	proxy *Proxy
}

func (lrp *LivereloadProxy) SetProxy(proxy config.Proxy, source config.Source) {
	lrp.proxy = NewProxy(proxy, source)
}

func (lrp *LivereloadProxy) AddTask(name string, taskConfig config.Task) {
	lrp.tasks = append(lrp.tasks, NewTask(name, lrp.proxy, taskConfig))
}

func (lrp *LivereloadProxy) Run() {
	lrp.runTasks()
	lrp.runProxy()
}

func (lrp *LivereloadProxy) runTasks() {
	for _, task := range lrp.tasks {
		task.Run()
	}
}

func (lrp *LivereloadProxy) runProxy() {
	lrp.proxy.Run()
}

func (lrp *LivereloadProxy) Stop() {
	lrp.stopProxy()
	lrp.stopTasks()
}

func (lrp *LivereloadProxy) stopTasks() {
	for _, task := range lrp.tasks {
		task.Stop()
	}
}

func (lrp *LivereloadProxy) stopProxy() {
	lrp.proxy.Close()
}
