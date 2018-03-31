package logger

import (
	"log"
)

func Info(name string, message string) {
	log.Printf("[%s] %s", name, message)
}

func InfoMonitor(name string, message string) {
	log.Printf("[%s] [monitor] %s", name, message)
}

func InfoCommand(name string, commandName string, message string) {
	log.Printf("[%s] [command] [%s] %s", name, commandName, message)
}

func InfoCommandStdout(name string, commandName string, message string) {
	log.Printf("[%s] [command] [%s] [stdout] %s", name, commandName, message)
}
