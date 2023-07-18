package ffmpeg

import (
	"IPTV_ReStreamer_GoLang/Logger"
	"bufio"
	"fmt"
	"github.com/opencoff/go-logger"
	"io"
	"os"
	"os/exec"
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
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.isFFmpegRunning {
		return fmt.Errorf("FFmpeg is already running")
	}
	ffmpegArgs := f.ConstructFFmpegArgs(args)

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

func (f *FFmpeg) processOutput(pipe io.Reader) {
	scanner := bufio.NewScanner(pipe)
	keywords := []string{
		"Opening",
		"Input",
		"keyword3",
	} // Add your desired keywords here

	for scanner.Scan() {
		line := scanner.Text()
		log.Info(line)
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
