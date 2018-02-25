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

func InfoCommand(name string, event string, commandName string, message string) {
	log.Printf("[%s] [%s] [%s] %s", name, commandName, event, message)
}
