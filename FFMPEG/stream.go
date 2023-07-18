package ffmpeg

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"reflect"
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
		InputFile: "input.mp4",
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
	bitRate1 := BitRate{Video: 1024, Audio: 128}
	bitRate2 := BitRate{Video: 2048, Audio: 256}

	resolution1 := Resolution{Width: 1920, Height: 1080}
	resolution2 := Resolution{Width: 1280, Height: 720}

	// создаем слайсы BitRate и Resolution
	bitRates := []BitRate{bitRate1, bitRate2}
	resolutions := []Resolution{resolution1, resolution2}

	// создаем экземпляр структуры Args и заполняем его
	args = Args{
		InputFile:  "http://romaxa55.otttv.pw/iptv/C2VHZLSGAWET4C/15117/index.m3u8",
		OutputDir:  "output",
		Resolution: resolutions,
		BitRate:    bitRates,
	}

	log.Debug("%+v\n", args)

	// Construct filter complex part
	filterComplex := f.filterComplex(args.Resolution)

	// Construct map video and audio part
	mapVideoAndAudio := f.mapVideoAndAudio(args.BitRate)

	varStreamMap := f.constructVarStreamMap(args.BitRate)

	ffmpegArgs := []string{
		"-i", args.InputFile,
		"-filter_complex", filterComplex,
	}

	// Append map video and audio arguments
	ffmpegArgs = append(ffmpegArgs, mapVideoAndAudio...)

	// Continue with the remaining arguments
	ffmpegArgs = append(ffmpegArgs,
		"-f", "hls",
		"-hls_time", "1",
		"-hls_list_size", "6",
		"-hls_flags", "delete_segments+omit_endlist+append_list",
		"-hls_playlist_type", "event",
		"-master_pl_name", "master.m3u8",
		"-hls_segment_filename", filepath.Join(args.OutputDir, "stream_%v", "data%d.ts"),
		"-var_stream_map", varStreamMap,
		filepath.Join(args.OutputDir, "stream_%v", "stream.m3u8"),
	)

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
