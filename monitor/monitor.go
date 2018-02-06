package monitor

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func NewMonitor(name string, task Task) *Monitor {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	return &Monitor{
		Name:    name,
		Task:    task,
		Watcher: watcher,
	}
}

type Monitor struct {
	Name    string
	Task    Task
	Watcher *fsnotify.Watcher
	Command *exec.Cmd
}

func (m *Monitor) Run(fn func(string)) {
	m.InitMonitorPath()
	m.StartCommand()

	// @todo log

	defer m.Watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-m.Watcher.Events:
				log.Info("event:", event)
				switch {
				case event.Op&fsnotify.Write == fsnotify.Write:
					log.Info("Modified file: ", event.Name)
				case event.Op&fsnotify.Create == fsnotify.Create:
					m.AddMonitorPathByString(event.Name)
					log.Info("Created file: ", event.Name)
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					log.Info("Removed file: ", event.Name)
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					log.Info("Renamed file: ", event.Name)
				case event.Op&fsnotify.Chmod == fsnotify.Chmod:
					log.Info("File changed permission: ", event.Name)
				}

				m.restartCommand()
				fn(event.Name)
			case err := <-m.Watcher.Errors:
				log.Info("error:", err)
			}
		}
	}()
	<-done
}

func (m *Monitor) InitMonitorPath() {
	wd := m.GetWorkingDirectory()

	err := filepath.Walk(wd, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		m.AddMonitorPath(path, fileInfo)
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (m *Monitor) AddMonitorPathByString(path string) bool {
	fileInfo, _ := os.Stat(path)
	return m.AddMonitorPath(path, fileInfo)
}

func (m *Monitor) AddMonitorPath(path string, fileInfo os.FileInfo) bool {
	if fileInfo == nil {
		return false
	}

	if !(fileInfo.IsDir()) {
		return false
	}

	err := m.Watcher.Add(path)
	if err != nil {
		panic(err)
	}

	return true
}

func (m *Monitor) GetWorkingDirectory() string {
	if m.Task.Path == "" {
		panic("task path is required.")
	}

	directory, err := filepath.Abs(m.Task.Path)
	if err != nil {
		panic(err)
	}

	return directory
}

func (m *Monitor) StartCommand() {
	m.Command = exec.Command("bash", "-c", m.Task.Command)
	stdout, _ := m.Command.StdoutPipe()
	m.Command.Start()
	oneByte := make([]byte, 100)
	go func() {
		for {
			_, err := stdout.Read(oneByte)
			if err != nil {
				log.Info(err.Error())
				break
			}
			r := bufio.NewReader(stdout)
			line, _, _ := r.ReadLine()
			log.Info(string(line))
		}
		m.Command.Wait()
	}()
}

func (m *Monitor) restartCommand() {
	m.KillCommand()
	m.StartCommand()
}

func (m *Monitor) KillCommand() {
	if m.Command.Process == nil {
		return
	}

	if !(m.Task.Compile) {
		return
	}

	log.Info("kill")
	m.Command.Process.Kill()
}
