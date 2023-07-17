package main

import (
	"IPTV_ReStreamer_GoLang/FFMPEG"
	"IPTV_ReStreamer_GoLang/Logger"
	"IPTV_ReStreamer_GoLang/Server"
	"fmt"
	"github.com/opencoff/go-logger"
)

var log *logger.Logger

func main() {
	err := Logger.InitLogger()
	if err != nil {
		fmt.Print(err)
		return
	}

	go func() {
		Server.StartServer()
	}()

	// Запуск процесса ffmpeg в отдельной горутине
	go func() {
		args := ffmpeg.NewArgs()
		fmt.Println(args)
		args.InputFile = "http://romaxa55.otttv.pw/iptv/C2VHZLSGAWET4C/15117/index.m3u8"

		err := ffmpeg.StartFFmpeg(args)
		if err != nil {
			fmt.Println("Error starting FFmpeg:", err)
			return
		}

		fmt.Println("FFmpeg command started successfully")
	}()

	// Бесконечный цикл для предотвращения завершения программы
	select {}
}
