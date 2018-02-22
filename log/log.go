package log

import (
	log "github.com/sirupsen/logrus"
)

func Info(name string, message string) {
	log.Info("[" + name + "] " + message)
}
