package ffmpeg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

func StartFFmpeg(args Args) error {
	ffmpegArgs := ConstructFfmpegArgs(args)

	cmd := exec.Command("ffmpeg", ffmpegArgs...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to start FFmpeg command: %w", err)
	}

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
