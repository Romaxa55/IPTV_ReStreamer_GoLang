package Server

import (
	"IPTV_ReStreamer_GoLang/Config"
	ffmpeg "IPTV_ReStreamer_GoLang/FFMPEG"
	"IPTV_ReStreamer_GoLang/Logger"
	"fmt"
	"github.com/opencoff/go-logger"
	"net/http"
	"os"
	"regexp"
	"strings"
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
}

func loggingMiddleware(next http.Handler, log *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("Request: %s %s %s", r.RemoteAddr, r.Method, r.URL))
		next.ServeHTTP(w, r)
	})
}

func StartServer() {
	f = ffmpeg.NewFFmpeg()
	args := ffmpeg.NewArgs()
	config := Config.GetServerConfig()
	log.Info(fmt.Sprintf("Starting server at http://%s:%s", config.IP, config.Port))

	// Новый таймер
	stopTimer := time.NewTimer(time.Minute * 1) // Например, одна минута

	// Запускаем горутину, которая будет ожидать срабатывание таймера
	go func() {
		<-stopTimer.C
		if err := f.StopFFmpeg(); err != nil {
			log.Error("Failed to stop FFmpeg:", err)
		}
	}()

	http.Handle("/iptv.m3u8", loggingMiddleware(http.StripPrefix("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("hash")
		if token != config.Token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		_, err := w.Write([]byte(playlist.Content))
		if err != nil {
			log.Error("Failed to write response:", err)
		}

	})), log))

	http.Handle("/", loggingMiddleware(http.StripPrefix("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stopTimer.Reset(time.Minute * 1)
		if matched, _ := regexp.MatchString(`^/intro_\d+\.ts$`, r.RequestURI); matched {
			startSegmentHandler(w, r)
			return
		}
		// Создаем регулярное выражение для поиска строки с доменом
		re := regexp.MustCompile(`.*` + regexp.QuoteMeta(r.RequestURI) + `.*\n`)

		// Ищем строку с доменом в плейлисте
		match := re.FindString(playlist.Original)

		// Если строка найдена, используем ее как входной файл
		if match != "" && args.InputFile != strings.TrimSpace(match) {
			args.InputFile = strings.TrimSpace(match)
			stoptStreamHandler(w, r)
			handlePlaylistRequest(w, r, args)
			return
		} else {
			handlePlaylistRequest(w, r, args)
			return
		}

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
