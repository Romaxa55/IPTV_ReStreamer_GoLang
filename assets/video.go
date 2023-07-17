package assets

import (
	"embed"
	"io"
	"io/fs"
)

//go:embed Video/*
var videoFolder embed.FS

// VideoFS returns the embedded video file system.
func VideoFS() fs.FS {
	return videoFolder
}

func UseVideoFS(fileName string) ([]byte, error) {
	videoFS := VideoFS()

	// Открываем файл
	file, err := videoFS.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Читаем файл
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Возвращаем содержимое файла
	return data, nil
}
