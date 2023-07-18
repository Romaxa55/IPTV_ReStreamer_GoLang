package ffmpeg

import (
	"IPTV_ReStreamer_GoLang/Logger"
	"bufio"
	"fmt"
	"github.com/opencoff/go-logger"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	LogApp = &Logger.App{}
	log    *logger.Logger
)

func init() {
	err := LogApp.InitLogger("ffmpeg") // Initialize the logger with "Server" as the prefix
	if err != nil {
		panic(err)
	}
	log = LogApp.Log // Set log equal to LogApp.Log
}

func (f *FFmpeg) StartFFmpeg(args Args) error {
	//ffmpegArgs := ConstructFfmpegArgs(args)
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.isFFmpegRunning {
		return fmt.Errorf("FFmpeg is already running")
	}

	ffmpegArgs := []string{
		"-i", args.InputFile,
		"-filter_complex",
		"[0:v]split=3[v1][v2][v3];" +
			"[v1]scale=w=1920:h=1080[v1out];" +
			"[v2]scale=w=1280:h=720[v2out];" +
			"[v3]scale=w=640:h=360[v3out]",
		"-map", "[v1out]", "-c:v:0", "libx265", "-b:v:0", "5000k",
		"-map", "[v2out]", "-c:v:1", "libx265", "-b:v:1", "3000k",
		"-map", "[v3out]", "-c:v:2", "libx265", "-b:v:2", "1000k",
		"-map", "a:0", "-c:a:0", "aac", "-b:a:0", "64k",
		"-map", "a:0", "-c:a:1", "aac", "-b:a:1", "64k",
		"-map", "a:0", "-c:a:2", "aac", "-b:a:2", "64k",
		"-f", "hls",
		"-hls_time", "10",
		"-hls_list_size", "20",
		"-hls_flags", "delete_segments+omit_endlist+append_list+program_date_time+discont_start",
		"-hls_playlist_type", "event",
		"-master_pl_name", "master.m3u8",
		"-hls_segment_filename", filepath.Join(args.OutputDir, "stream_%v", "data%d.ts"),
		"-var_stream_map", "v:0,a:0 v:1,a:1 v:2,a:2",
		filepath.Join(args.OutputDir, "stream_%v", "stream.m3u8"),
	}
	f.cmd = exec.Command("ffmpeg", ffmpegArgs...)
	stderrPipe, err := f.cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdoutPipe, err := f.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = f.cmd.Start()
	if err != nil {
		return err
	}
	go f.processOutput(stderrPipe)
	go f.processOutput(stdoutPipe)

	f.isFFmpegRunning = true

	go func() {
		err := f.cmd.Wait()
		f.mutex.Lock()
		if err != nil {
			log.Error("FFmpeg process exited with error:", err)
		}
		f.isFFmpegRunning = false
		f.mutex.Unlock()
	}()

	return nil
}

func (f *FFmpeg) StopFFmpeg() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if !f.isFFmpegRunning {
		return nil
	}

	err := f.cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return err
	}

	f.isFFmpegRunning = false
	return nil
}

func ConstructFfmpegArgs(args Args) []string {
	ffmpegArgs := []string{
		"-i", args.InputFile,
		"-filter_complex",
		fmt.Sprintf("[0:v]split=3[v1][v2][v3];"),
	}

	videoMapArgs := generateVideoMapArgs(args.VideoScaleResolution, args.VideoCodecBitrate)
	audioMapArgs := generateAudioMapArgs(args.AudioCodecBitrate)

	ffmpegArgs = append(ffmpegArgs, videoMapArgs...)
	ffmpegArgs = append(ffmpegArgs, audioMapArgs...)

	ffmpegArgs = append(ffmpegArgs,
		"-f", "hls",
		"-hls_time", args.HlsTime,
		"-hls_list_size", args.HlsListSize,
		"-hls_flags", args.HlsFlags,
		"-hls_playlist_type", args.HlsPlaylistType,
		"-master_pl_name", args.MasterPlaylistName,
		"-hls_segment_filename", filepath.Join(args.OutputDir, args.OutputTemplate),
		"-var_stream_map", args.AudioStreamMap,
		filepath.Join(args.OutputDir, args.OutputTemplate),
	)

	return ffmpegArgs
}

func generateVideoFilter(scales []string) string {
	filter := ""
	for i, scale := range scales {
		filter += fmt.Sprintf("[v%d]scale=w=%s:h=%s[v%dout];", i, scale, scale, i)
	}
	return filter
}

func generateVideoMapArgs(scales []string, bitrates []string) []string {
	var args []string
	for i := 0; i < len(scales); i++ {
		args = append(args, "-map", fmt.Sprintf("[v%dout]", i), "-c:v:"+strconv.Itoa(i), "libx265", "-b:v:"+strconv.Itoa(i), bitrates[i])
	}
	return args
}

func generateAudioMapArgs(bitrates []string) []string {
	var args []string
	for i := 0; i < len(bitrates); i++ {
		args = append(args, "-map", "a:0", "-c:a:"+strconv.Itoa(i), "aac", "-b:a:"+strconv.Itoa(i), bitrates[i])
	}
	return args
}

func (f *FFmpeg) processOutput(pipe io.Reader) {
	scanner := bufio.NewScanner(pipe)
	keywords := []string{
		"Opening",
		"Input",
		"keyword3",
	} // Add your desired keywords here

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "error") {
			log.Error(line)
		} else {
			for _, keyword := range keywords {
				if strings.Contains(line, keyword) {
					log.Info(line)
					break // Break the loop once a keyword match is found
				}
			}
		}
	}
}
