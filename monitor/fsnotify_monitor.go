package monitor

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/oneut/lrp/config"
	"github.com/oneut/lrp/logger"
)

func NewFsnotifyMonitor(name string, monitorConfig config.Monitor) *fsnotifyMonitor {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	return &fsnotifyMonitor{
		Name:          name,
		MonitorConfig: monitorConfig,
		Watcher:       watcher,
	}
}

type fsnotifyMonitor struct {
	Name          string
	MonitorConfig config.Monitor
	Watcher       *fsnotify.Watcher
}

func (fm *fsnotifyMonitor) Run(fn func(string)) {
	logger.InfoMonitor(fm.Name, "start")
	defer fm.Watcher.Close()
	fm.initMonitorPath()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-fm.Watcher.Events:
				absPath := fm.getAbsPath(event.Name)
				if fm.isIgnoreAbsPath(absPath) {
					break
				}
				switch {
				case event.Op&fsnotify.Create == fsnotify.Create:
					fm.addMonitorPath(absPath)
					fn(absPath)
					logger.InfoMonitor(fm.Name, event.Op.String()+" "+absPath)
				case event.Op&fsnotify.Write == fsnotify.Write:
					fm.addMonitorPath(absPath)
					fn(absPath)
					logger.InfoMonitor(fm.Name, event.Op.String()+" "+absPath)
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					fm.addMonitorPath(absPath)
					fn(absPath)
					logger.InfoMonitor(fm.Name, event.Op.String()+" "+absPath)
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					// `remove monitor path` is not required
					fn(absPath)
					logger.InfoMonitor(fm.Name, event.Op.String()+" "+absPath)
				}

			case err := <-fm.Watcher.Errors:
				panic(err)
			}
		}
	}()
	<-done
}

func (fm *fsnotifyMonitor) isIgnoreAbsPath(absPath string) bool {
	if len(fm.MonitorConfig.Ignore) == 0 {
		return false
	}

	for _, ignore := range fm.MonitorConfig.Ignore {
		if strings.Contains(absPath, ignore) {
			return true
		}

		absIgnore := fm.getAbsPath(ignore)
		if strings.HasPrefix(absPath, absIgnore) {
			return true
		}

		absIgnoreRegexp, err := regexp.Compile(absIgnore)
		if err != nil {
			panic(err)
		}

		if absIgnoreRegexp.MatchString(absPath) {
			return true
		}
	}
	return false
}

func (fm *fsnotifyMonitor) initMonitorPath() error {
	for _, targetPath := range fm.MonitorConfig.Paths {
		absTargetPath := fm.getAbsPath(targetPath)
		if absTargetPath == "" {
			return nil
		}
		err := filepath.Walk(absTargetPath, func(path string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if fm.isIgnoreAbsPath(path) {
				return nil
			}

			fm.addMonitorPath(path)
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (fm *fsnotifyMonitor) addMonitorPath(path string) bool {
	var err error
	var fileInfo os.FileInfo
	fileInfo, err = os.Stat(path)
	if err != nil {
		return false
	}

	if !(fileInfo.IsDir()) {
		return false
	}

	err = fm.Watcher.Add(path)
	if err != nil {
		panic(err)
	}

	logger.InfoMonitor(fm.Name, "watch "+path)
	return true
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
	logger.InfoMonitor(fm.Name, "stop")
	fm.Watcher.Close()
}
