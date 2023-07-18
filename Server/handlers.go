package Server

import (
	"IPTV_ReStreamer_GoLang/FFMPEG"
	"fmt"
	"github.com/AlekSi/pointer"
	"github.com/etherlabsio/go-m3u8/m3u8"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func handlePlaylistRequest(w http.ResponseWriter, r *http.Request, args ffmpeg.Args) {
	fileServer := http.FileServer(http.Dir("output"))

	if strings.Contains(r.RequestURI, "m3u8") {
		startStreamHandler(w, r, args)

		// Проверяем, существует ли файл master.m3u8
		if _, err := os.Stat("output/master.m3u8"); err == nil {
			// Файл существует, читаем его и возвращаем содержимое
			data, err := os.ReadFile("output/master.m3u8")
			if err != nil {
				log.Error("Failed to read master.m3u8:", err)
				return
			}
			// Регулярное выражение для поиска путей, оканчивающихся на .m3u8
			regex := regexp.MustCompile(`(?m)^(.*\.m3u8)$`)

			// Заменить пути, добавив / в начале
			modifiedData := regex.ReplaceAllString(string(data), "/$1")

			// Устанавливаем заголовок и возвращаем содержимое файла
			_, err = w.Write([]byte(modifiedData))
			if err != nil {
				log.Error("Failed to write data:", err)
			}
		} else if os.IsNotExist(err) {
			playlist := &m3u8.Playlist{
				Version:             pointer.ToInt(3),
				Cache:               pointer.ToBool(false),
				Live:                false,
				IndependentSegments: true,
				Items: []m3u8.Item{
					&m3u8.SegmentItem{
						Duration: 10,
						Segment:  "/intro_00000.ts",
					},
					&m3u8.SegmentItem{
						Duration: 10,
						Segment:  "/intro_00001.ts",
					},
					&m3u8.SegmentItem{
						Duration: 10,
						Segment:  "/intro_00002.ts",
					},
					&m3u8.SegmentItem{
						Duration: 10,
						Segment:  "/intro_00003.ts",
					},
				},
			}

			w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(playlist.String())))
			w.Write([]byte(playlist.String()))

		} else {
			// Неизвестная ошибка
			log.Error("Failed to check master.m3u8:", err)
		}
	}

	//pathPattern := regexp.MustCompile(`m3u8`)
	//if pathPattern.MatchString(r.URL.Path) {
	//	// Обработка запроса для /stream_2/stream.m3u8
	//	// Можно вызвать startStreamHandler или выполнить другие нужные действия
	//	startStreamHandler(w, r, args)
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
