package Server

import (
	"IPTV_ReStreamer_GoLang/FFMPEG"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func handlePlaylistRequest(w http.ResponseWriter, r *http.Request, args ffmpeg.Args) {

	if strings.Contains(r.RequestURI, "m3u8") && !strings.Contains(r.RequestURI, "stream.m3u8") {
		startStreamHandler(w, r, args)

		// Проверяем, существует ли файл master.m3u8
		masterFile := filepath.Join(args.OutputDir, "master.m3u8")

		if _, err := os.Stat(masterFile); err == nil {
			// Файл существует, читаем его и возвращаем содержимое
			data, err := os.ReadFile(masterFile)
			if err != nil {
				log.Error("Failed to read master.m3u8:", err)
				return
			}
			// Регулярное выражение для поиска путей, оканчивающихся на .m3u8
			regex := regexp.MustCompile(`(?m)^(.*\.m3u8)$`)

			// Заменить пути, добавив / в начале
			modifiedData := regex.ReplaceAllString(string(data), "/$1")

			// Проверка существования файлов "data*.ts" в каждом stream
			streams := regex.FindAllString(string(data), -1)
			success := false
			for _, stream := range streams {
				directory := filepath.Dir(stream)
				directory = filepath.Join(args.OutputDir, directory)
				i := 0
				for {
					filename := fmt.Sprintf("%s/data%d.ts", directory, i)
					if _, err := os.Stat(filename); os.IsNotExist(err) {
						break // файл не существует, прекращаем поиск
					} else if err != nil {
						log.Error("Failed to stat ", filename, ":", err)
						return // произошла ошибка, прекращаем выполнение
					}
					i++
					if i > 1 { // Успех, если найдено более одного файла
						success = true
						break
					}
				}
				if success {
					break
				}
			}

			if success {
				// Устанавливаем заголовок и возвращаем содержимое файла
				_, err = w.Write([]byte(modifiedData))
				if err != nil {
					log.Error("Failed to write data:", err)
				}
			} else {
				sendPlaylist(w)
			}
		} else if os.IsNotExist(err) {
			sendPlaylist(w)
		} else {
			// Неизвестная ошибка
			log.Error("Failed to check master.m3u8:", err)
		}
	}

	pathPattern := regexp.MustCompile(`stream.m3u8`)
	if pathPattern.MatchString(r.URL.Path) {
		// Получить путь до каталога stream из URL
		directory := filepath.Dir(r.URL.Path)

		// Путь до файла stream.m3u8
		streamFile := filepath.Join(args.OutputDir, directory, "stream.m3u8")
		// Проверяем, существует ли файл stream.m3u8
		if _, err := os.Stat(streamFile); err == nil {
			// Файл существует, читаем его и возвращаем содержимое
			data, err := os.ReadFile(streamFile)
			if err != nil {
				log.Error("Failed to read stream.m3u8:", err)
				return
			}
			// Регулярное выражение для поиска путей, оканчивающихся на .ts
			regex := regexp.MustCompile(`(?m)^(.*\.ts)$`)

			// Заменить пути, добавив /directory/ в начале
			modifiedData := regex.ReplaceAllStringFunc(string(data), func(s string) string {
				return "/" + directory + "/" + s
			})

			// Здесь вы можете добавить дополнительную обработку данных из stream.m3u8, как вы делали с master.m3u8

			// Устанавливаем заголовок и возвращаем содержимое файла
			_, err = w.Write([]byte(modifiedData))
			if err != nil {
				log.Error("Failed to write data:", err)
			}
			return
		} else {
			// Выводим ошибку, если файл не найден
			sendPlaylist(w)
		}

		fileServer := http.FileServer(http.Dir(args.OutputDir))
		fileServer.ServeHTTP(w, r)
	}

	pathPattern = regexp.MustCompile(`.ts`)
	if pathPattern.MatchString(r.URL.Path) {
		fileServer := http.FileServer(http.Dir(args.OutputDir))
		fileServer.ServeHTTP(w, r)
	}

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
