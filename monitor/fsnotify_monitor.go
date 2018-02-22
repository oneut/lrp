package monitor

import (
	"github.com/fsnotify/fsnotify"
	"github.com/oneut/lrp/config"
	"github.com/oneut/lrp/log"
	"os"
	"path/filepath"
)

func NewFsnotifyMonitor(name string, monitorConfig config.Monitor) *FsnotifyMonitor {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	return &FsnotifyMonitor{
		Name:          name,
		MonitorConfig: monitorConfig,
		Watcher:       watcher,
	}
}

type FsnotifyMonitor struct {
	Name          string
	MonitorConfig config.Monitor
	Watcher       *fsnotify.Watcher
}

func (fm *FsnotifyMonitor) Run(fn func(string)) {
	defer fm.Watcher.Close()
	fm.InitMonitorPath()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-fm.Watcher.Events:
				/*		if event.Name[len(event.Name)-1:] == "~" {
							continue
						}
						if event.Name[len(event.Name)-4:] == ".swp" {
							continue
						}
				*/log.Info(fm.Name, "event:"+event.Name)
				switch {
				case event.Op&fsnotify.Create == fsnotify.Create:
					fm.AddMonitorPathByString(event.Name)
					fn(event.Name)
				case event.Op&fsnotify.Create == fsnotify.Write:
					fm.AddMonitorPathByString(event.Name)
					fn(event.Name)
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					fm.AddMonitorPathByString(event.Name)
					fn(event.Name)
				}

			case err := <-fm.Watcher.Errors:
				panic(err)
			}
		}
	}()
	<-done
}

func (fm *FsnotifyMonitor) InitMonitorPath() error {
	for _, targetPath := range fm.MonitorConfig.Paths {
		absTargetPath := fm.GetAbsPath(targetPath)
		if absTargetPath == "" {
			return nil
		}
		err := filepath.Walk(absTargetPath, func(path string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fm.AddMonitorPath(path, fileInfo)
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (fm *FsnotifyMonitor) AddMonitorPathByString(path string) bool {
	fileInfo, _ := os.Stat(path)
	return fm.AddMonitorPath(path, fileInfo)
}

func (fm *FsnotifyMonitor) AddMonitorPath(path string, fileInfo os.FileInfo) bool {
	if fileInfo == nil {
		return false
	}

	if !(fileInfo.IsDir()) {
		return false
	}

	err := fm.Watcher.Add(path)
	if err != nil {
		panic(err)
	}

	return true
}

func (fm *FsnotifyMonitor) GetAbsPath(path string) string {
	if path == "" {
		return ""
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return absPath
}
