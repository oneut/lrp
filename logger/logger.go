package logger

import (
	"log"
)

func Info(name string, message string) {
	log.Printf("[%s] %s", name, message)
}

func InfoEvent(name string, event string, message string) {
	log.Printf("[%s] [%s] %s", name, event, message)
}

func InfoCommand(name string, commandName string, message string) {
	log.Printf("[%s] [command] [%s] %s", name, commandName, message)
}

func InfoCommandStdout(name string, commandName string, message string) {
	log.Printf("[%s] [command] [%s] [stdout] %s", name, commandName, message)
}
