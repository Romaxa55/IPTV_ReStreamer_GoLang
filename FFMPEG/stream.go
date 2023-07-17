package ffmpeg

import "reflect"

type Args struct {
	InputFile            string
	OutputDir            string
	OutputTemplate       string
	VideoScaleResolution []string
	VideoCodecBitrate    []string
	AudioCodecBitrate    []string
	AudioStreamMap       string
	HlsTime              string
	HlsListSize          string
	HlsFlags             string
	HlsPlaylistType      string
	MasterPlaylistName   string
}

func NewArgs() Args {
	args := Args{}

	// Задаем значения по умолчанию только для неустановленных полей
	defaults := Args{
		InputFile:            "input.mp4",
		OutputDir:            "output",
		OutputTemplate:       "stream_%v/stream.m3u8",
		VideoScaleResolution: []string{"1920x1080", "1280x720", "640x360"},
		VideoCodecBitrate:    []string{"5000k", "3000k", "1000k"},
		AudioCodecBitrate:    []string{"64k", "64k", "64k"},
		AudioStreamMap:       "v:0,a:0 v:1,a:1 v:2,a:2",
		HlsTime:              "10",
		HlsListSize:          "20",
		HlsFlags:             "delete_segments+omit_endlist+append_list+program_date_time+discont_start",
		HlsPlaylistType:      "event",
		MasterPlaylistName:   "master.m3u8",
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
