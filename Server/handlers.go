package Server

import (
	"IPTV_ReStreamer_GoLang/FFMPEG"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func handlePlaylistRequest(w http.ResponseWriter, r *http.Request, args ffmpeg.Args) {
	fileServer := http.FileServer(http.Dir("output"))

	if strings.Contains(r.RequestURI, "master.m3u8") {
		startStreamHandler(w, r, args)

		// Проверяем, существует ли файл master.m3u8
		if _, err := os.Stat("1output/master.m3u8"); err == nil {
			// Файл существует, читаем его и возвращаем содержимое
			data, err := os.ReadFile("output/master.m3u8")
			if err != nil {
				log.Error("Failed to read master.m3u8:", err)
				return
			}

			// Устанавливаем заголовок и возвращаем содержимое файла
			w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
			_, err = w.Write(data)
			if err != nil {
				log.Error("Failed to write data:", err)
			}
		} else if os.IsNotExist(err) {
			playlist := &m3u8.MediaPlaylist{
				TargetDuration:   10,
				MediaSequence:    0,
				PlaylistType:     m3u8.EVENT,
				DiscontinuitySeq: 0,
				Version:          3,
				Segments: []*m3u8.MediaSegment{
					{
						Duration: 10,
						URI:      "intro_00000.ts",
					},
					{
						Duration: 10,
						URI:      "intro_00001.ts",
					},
					// Добавьте другие сегменты в соответствии с вашими требованиями
				},
			}
			w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(playlistContent)))
			w.Write([]byte(playlistContent))

		} else {
			// Неизвестная ошибка
			log.Error("Failed to check master.m3u8:", err)
		}
	}

	pathPattern := regexp.MustCompile(`m3u8`)
	if pathPattern.MatchString(r.URL.Path) {
		// Обработка запроса для /stream_2/stream.m3u8
		// Можно вызвать startStreamHandler или выполнить другие нужные действия
		startStreamHandler(w, r, args)
	}

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
