package ffmpeg

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type Args struct {
	InputFile  string
	OutputDir  string
	Resolution []Resolution
	BitRate    []BitRate
}

type Resolution struct {
	Width  int
	Height int
}

type BitRate struct {
	Video int
	Audio int
}

type FFmpeg struct {
	mutex           sync.Mutex
	cmd             *exec.Cmd
	isFFmpegRunning bool
}

func NewFFmpeg() *FFmpeg {
	return &FFmpeg{
		cmd:             nil,
		isFFmpegRunning: false,
	}
}

func NewArgs() Args {
	args := Args{}

	// Задаем значения по умолчанию только для неустановленных полей
	defaults := Args{
		InputFile: "intro_00000.ts",
		OutputDir: "output",
	}

	// Получаем указатели на поля структуры
	argsValue := reflect.ValueOf(&args).Elem()
	defaultsValue := reflect.ValueOf(defaults)

	// Проверяем каждое поле и устанавливаем значение по умолчанию, если поле не было установлено
	for i := 0; i < argsValue.NumField(); i++ {
		fieldValue := argsValue.Field(i)
		defaultValue := defaultsValue.Field(i)
		if fieldValue.IsZero() {
			fieldValue.Set(defaultValue)
		}
	}

	return args
}

func (f *FFmpeg) ConstructFFmpegArgs(args Args) []string {
	// создаем экземпляры структур BitRate и Resolution
	//config := Config.GetServerConfig()
	bitRate1 := BitRate{Video: 64, Audio: 32}

	ffmpegArgs := []string{
		"-i", args.InputFile,
		"-c:v", "libx265", // Используйте кодек H.265 для видео
		"-b:v", strconv.Itoa(bitRate1.Video) + "k", // Установите битрейт для видео
		"-b:a", strconv.Itoa(bitRate1.Audio) + "k", // Установите битрейт для аудио
		"-f", "hls",
		"-hls_time", "2",
		"-hls_list_size", "11",
		"-hls_flags", "delete_segments",
		"-master_pl_name", "master.m3u8",
		"-hls_playlist_type", "event",
		"-hls_segment_filename", filepath.Join(args.OutputDir, "stream_%v", "data%d.ts"),
		filepath.Join(args.OutputDir, "stream_%v", "stream.m3u8"),
	}

	return ffmpegArgs
}

func (f *FFmpeg) filterComplex(res []Resolution) string {
	split := fmt.Sprintf("[0:v]split=%d", len(res))

	var names []string
	for i := range res {
		names = append(names, fmt.Sprintf("[v%d]", i+1))
	}

	var scales []string
	for i, r := range res {
		scales = append(scales, fmt.Sprintf("[v%d]scale=w=%d:h=%d[v%dout]", i+1, r.Width, r.Height, i+1))
	}

	return split + strings.Join(names, "") + ";" + strings.Join(scales, ";")
}

func (f *FFmpeg) mapVideoAndAudio(br []BitRate) []string {
	var result []string
	for i, b := range br {
		videoOpts := fmt.Sprintf("-map [v%dout] -c:v:%d libx265 -b:v:%d %dk", i+1, i, i, b.Video)
		audioOpts := fmt.Sprintf("-map a:0 -c:a:%d aac -b:a:%d %dk", i, i, b.Audio)

		result = append(result, strings.Fields(videoOpts)...)
		result = append(result, strings.Fields(audioOpts)...)

	}

	return result
}

func (f *FFmpeg) constructVarStreamMap(br []BitRate) string {
	var result []string
	for i := range br {
		result = append(result, fmt.Sprintf("v:%d,a:%d", i, i))
	}
	return strings.Join(result, " ")
}
