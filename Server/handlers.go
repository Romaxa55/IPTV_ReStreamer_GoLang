package Server

import (
	"IPTV_ReStreamer_GoLang/FFMPEG"
	"net/http"
	"regexp"
)

func handlePlaylistRequest(w http.ResponseWriter, r *http.Request, args ffmpeg.Args) {
	fileServer := http.FileServer(http.Dir("output"))

	pathPattern := regexp.MustCompile(`m3u8`)
	if pathPattern.MatchString(r.URL.Path) {
		// Обработка запроса для /stream_2/stream.m3u8
		// Можно вызвать startStreamHandler или выполнить другие нужные действия
		startStreamHandler(w, r, args)
	}

	//// Обработка других типов запросов
	//if r.URL.Path == "/master.m3u" {
	//	// Если запрашивается master.m3u и его нет, можно временно предоставить альтернативный плейлист
	//	alternativePlaylist := "#EXTM3U\n#EXT-X-STREAM-INF:BANDWIDTH=800000,RESOLUTION=640x360\nalternative_stream.m3u8"
	//	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	//	_, err := w.Write([]byte(alternativePlaylist))
	//	if err != nil {
	//		log.Error("Failed to write alternative playlist:", err)
	//	}
	//
	//}

	// Обработка других запросов
	fileServer.ServeHTTP(w, r)
}

func startStreamHandler(w http.ResponseWriter, r *http.Request, args ffmpeg.Args) {
	err := f.StartFFmpeg(args)
	if err != nil {
		log.Error(err.Error())
	}
}

func stoptStreamHandler(w http.ResponseWriter, r *http.Request) {
	err := f.StopFFmpeg()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
