package monitor

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/oneut/lrp/config"
	"github.com/oneut/lrp/logger"
)

func NewFsnotifyMonitor(name string, monitorConfig config.Monitor) *fsnotifyMonitor {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	monitor := &fsnotifyMonitor{
		name:    name,
		watcher: watcher,
		paths:   monitorConfig.Paths,
	}

	for _, ignore := range monitorConfig.Ignores {
		if ignore != "" {
			monitor.AddIgnore(ignore)
		}
	}

	return monitor
}

type fsnotifyMonitor struct {
	name    string
	watcher *fsnotify.Watcher
	paths   []string
	ignores []*Ignore
}

func (fm *fsnotifyMonitor) AddIgnore(ignore string) {
	fm.ignores = append(fm.ignores, NewIgnore(ignore))
}

func (fm *fsnotifyMonitor) Run(fn func(string)) {
	logger.InfoMonitor(fm.name, "start")
	defer fm.watcher.Close()
	fm.initMonitorPath()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-fm.watcher.Events:
				absPath := fm.getAbsPath(event.Name)
				if fm.isIgnoreAbsPath(absPath) {
					break
				}
				switch {
				case event.Op&fsnotify.Create == fsnotify.Create:
					fm.addWalkMonitorPath(absPath)
					fn(absPath)
					logger.InfoMonitor(fm.name, event.Op.String()+" "+absPath)
				case event.Op&fsnotify.Write == fsnotify.Write:
					fm.addWalkMonitorPath(absPath)
					fn(absPath)
					logger.InfoMonitor(fm.name, event.Op.String()+" "+absPath)
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					fm.addWalkMonitorPath(absPath)
					fn(absPath)
					logger.InfoMonitor(fm.name, event.Op.String()+" "+absPath)
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					// `remove monitor path` is not required
					fn(absPath)
					logger.InfoMonitor(fm.name, event.Op.String()+" "+absPath)
				}

			case err := <-fm.watcher.Errors:
				panic(err)
			}
		}
	}()
	<-done
}

func (fm *fsnotifyMonitor) isIgnoreAbsPath(absPath string) bool {
	if len(fm.ignores) == 0 {
		return false
	}

	for _, ignore := range fm.ignores {
		if ignore.Match(absPath) {
			return true
		}
	}
	return false
}

func (fm *fsnotifyMonitor) initMonitorPath() {
	for _, path := range fm.paths {
		absPath := fm.getAbsPath(path)
		if absPath == "" {
			continue
		}

		fm.addWalkMonitorPath(absPath)
	}
}

func (fm *fsnotifyMonitor) addWalkMonitorPath(absPath string) {
	err := filepath.Walk(absPath, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			// not exists file or directory
			return nil
		}

		if fm.isIgnoreAbsPath(path) {
			return nil
		}

		if fileInfo.IsDir() {
			fm.addMonitorPath(path)
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (fm *fsnotifyMonitor) addMonitorPath(path string) {
	err := fm.watcher.Add(path)
	if err != nil {
		panic(err)
	}

	logger.InfoMonitor(fm.name, "watch "+path)
}

func (fm *fsnotifyMonitor) getAbsPath(path string) string {
	if path == "" {
		return ""
	}

	if filepath.IsAbs(path) {
		return path
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return absPath
}

func (fm *fsnotifyMonitor) Stop() {
	logger.InfoMonitor(fm.name, "stop")
	fm.watcher.Close()
}
