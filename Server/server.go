package Server

import (
	"IPTV_ReStreamer_GoLang/Config"
	ffmpeg "IPTV_ReStreamer_GoLang/FFMPEG"
	"IPTV_ReStreamer_GoLang/Logger"
	"fmt"
	"github.com/opencoff/go-logger"
	"net/http"
	"time"
)

var (
	f      *ffmpeg.FFmpeg
	LogApp = &Logger.App{}
	log    *logger.Logger
)

func init() {
	err := LogApp.InitLogger("Server") // Initialize the logger with "Server" as the prefix
	if err != nil {
		panic(err)
	}
	log = LogApp.Log // Set log equal to LogApp.Log
}

func loggingMiddleware(next http.Handler, log *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("Request: %s %s %s", r.RemoteAddr, r.Method, r.URL))
		next.ServeHTTP(w, r)
	})
}

func StartServer(args ffmpeg.Args) {
	f = ffmpeg.NewFFmpeg()
	config := Config.GetServerConfig()
	log.Info(fmt.Sprintf("Starting server at http://%s:%s", config.IP, config.Port))

	// Новый таймер
	stopTimer := time.NewTimer(time.Minute * 1) // Например, одна минуту

	// Запускаем горутину, которая будет ожидать срабатывание таймера
	go func() {
		<-stopTimer.C
		if err := f.StopFFmpeg(); err != nil {
			log.Error("Failed to stop FFmpeg:", err)
		}
	}()

	fileServer := http.FileServer(http.Dir("output"))
	http.Handle("/", loggingMiddleware(http.StripPrefix("/", fileServer), log))

	http.Handle("/start_stream", loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Сбрасываем таймер каждый раз при получении нового запроса
		stopTimer.Reset(time.Minute * 1)
		startStreamHandler(w, r, args)
	}), log))

	http.Handle("/stop_stream", loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Сбрасываем таймер каждый раз при получении нового запроса
		stoptStreamHandler(w, r)
	}), log))

	err := http.ListenAndServe(config.IP+":"+config.Port, nil)
	if err != nil {
		log.Error("Server failed to start: ", err)
	}
}

func startStreamHandler(w http.ResponseWriter, r *http.Request, args ffmpeg.Args) {
	err := f.StartFFmpeg(args)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func stoptStreamHandler(w http.ResponseWriter, r *http.Request) {
	err := f.StopFFmpeg()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
