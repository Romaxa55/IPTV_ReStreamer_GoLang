package Server

import (
	"IPTV_ReStreamer_GoLang/Config"
	ffmpeg "IPTV_ReStreamer_GoLang/FFMPEG"
	"IPTV_ReStreamer_GoLang/Logger"
	"fmt"
	"github.com/etherlabsio/go-m3u8"
	"github.com/opencoff/go-logger"
	"net/http"
	"os"
	"regexp"
	"time"
)

var (
	f      *ffmpeg.FFmpeg
	LogApp = &Logger.App{}
	log    *logger.Logger
	args   *ffmpeg.Args
)

func init() {
	err := LogApp.InitLogger("Server") // Initialize the logger with "Server" as the prefix
	if err != nil {
		panic(err)
	}
	log = LogApp.Log // Set log equal to LogApp.Log
	args := ffmpeg.NewArgs()
	args.InputFile = "http://romaxa55.otttv.pw/iptv/C2VHZLSGAWET4C/15117/index.m3u8"
}

func loggingMiddleware(next http.Handler, log *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("Request: %s %s %s", r.RemoteAddr, r.Method, r.URL))
		next.ServeHTTP(w, r)
	})
}

func StartServer() {
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

	http.Handle("/", loggingMiddleware(http.StripPrefix("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stopTimer.Reset(time.Minute * 1)
		if matched, _ := regexp.MatchString(`^/intro_\d+\.ts$`, r.RequestURI); matched {
			startSegmentHandler(w, r)
			return
		}
		handlePlaylistRequest(w, r, ffmpeg.Args{})

	})), log))

	err := http.ListenAndServe(config.IP+":"+config.Port, nil)
	if err != nil {
		log.Error("Server failed to start: ", err)
	}
}

func startSegmentHandler(w http.ResponseWriter, r *http.Request) {
	mainFileContent, err := os.ReadFile("Video/intro_00000.ts")
	if err != nil {
		http.Error(w, "Failed to read main file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "video/mp2t")
	_, err = w.Write(mainFileContent)
	if err != nil {
		log.Error("Failed to write response:", err)
	}
}
