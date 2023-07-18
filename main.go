package main

import (
	ffmpeg "IPTV_ReStreamer_GoLang/FFMPEG"
	"IPTV_ReStreamer_GoLang/Logger"
	"IPTV_ReStreamer_GoLang/Server"
	"fmt"
	"github.com/opencoff/go-logger"
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

	go func() {
		args := ffmpeg.NewArgs()
		args.InputFile = "http://romaxa55.otttv.pw/iptv/C2VHZLSGAWET4C/15117/index.m3u8"
		Server.StartServer(args)
	}()

	select {}
}
