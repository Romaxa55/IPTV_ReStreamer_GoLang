package main

import (
	"IPTV_ReStreamer_GoLang/Logger"
	"IPTV_ReStreamer_GoLang/Server"
	"fmt"
	"github.com/opencoff/go-logger"
	"os"
)

var (
	LogApp = &Logger.App{}
	log    *logger.Logger
)

func main() {
	err := LogApp.InitLogger("Main")
	if err != nil {
		fmt.Print(err)
		return
	}
	log = LogApp.Log        // Set log equal to LogApp.Log
	log.Info("Started app") // Use f.Log for logging

	// Удаляем папку "output"
	log.Info("Clean output") // Use f.Log for logging
	err = os.RemoveAll("output")
	if err != nil {
		log.Error("Failed to remove 'output' directory:", err)
	}
	go func() {
		Server.StartServer()
	}()

	select {}
}
