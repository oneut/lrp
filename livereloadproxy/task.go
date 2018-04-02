package livereloadproxy

import (
	"time"

	"github.com/oneut/lrp/command"
	"github.com/oneut/lrp/config"
	"github.com/oneut/lrp/monitor"
)

func NewTask(name string, proxy *Proxy, taskConfig config.Task) *Task {
	task := &Task{
		name:             name,
		proxy:            proxy,
		aggregateTimeout: taskConfig.GetAggregateTimeout(),
	}
	for commandName, commandConfig := range taskConfig.Commands {
		task.AddCommand(commandName, commandConfig)
	}

	task.SetMonitor(taskConfig.Monitor)
	return task
}

type Task struct {
	name             string
	proxy            *Proxy
	commands         []command.CommandInterface
	monitor          monitor.MonitorInterface
	isReloading      bool
	aggregateTimeout time.Duration
}

func (t *Task) AddCommand(commandName string, commandConfig config.Command) {
	t.commands = append(t.commands, command.NewCommand(t.name, commandName, commandConfig))
}

func (t *Task) SetMonitor(monitorConfig config.Monitor) {
	t.monitor = monitor.NewMonitor(t.name, monitorConfig)
}

func (t *Task) Run() {
	t.isReloading = false
	for _, cmd := range t.commands {
		go cmd.Run(t.Callback)
	}
	go t.monitor.Run(t.Callback)
}

func (t *Task) Callback(message string) {
	if t.isReloading {
		return
	}

	t.isReloading = true
	go func() {
		time.Sleep(t.aggregateTimeout)
		for _, cmd := range t.commands {
			if cmd.NeedsRestart() {
				cmd.Kill()
				go cmd.Start()
			}
		}
		t.proxy.Reload(message)
		t.isReloading = false
	}()
}

func (t *Task) Stop() {
	for _, cmd := range t.commands {
		cmd.Stop()
	}
	t.monitor.Stop()
}
