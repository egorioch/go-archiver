package service

import (
	"bytes"
	"fmt"
	"go-archiver/package/utils"
	"log"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

var (
	ffmpegRecordingIsAlreadyInProgress = "recording_is_already_in_progress"
	ffmpegFailedStartRecording         = "ffmpeg_failed_start_recording"
	ffmpegNoRecordingInProgress        = "ffmpeg_no_recording_in_progress"
	ffmpegErrorStopRecording           = "ffmpeg_error_stop_recording"
	ffmpegErrorWaiting                 = "ffmpeg_error_waiting"
)

type VideoRecorder struct {
	ID                    int
	mu                    sync.Mutex
	isRecording           bool
	cmd                   *exec.Cmd
	lastVideoFile         string
	logger                *log.Logger
	staticVideosDirectory string
	thumbnailsPath        string
	streamURL             string
}

func NewVideoRecorder(id int, logger *log.Logger, staticVideosDirectory, thumbnailsPath, streamURL string) *VideoRecorder {
	return &VideoRecorder{
		ID:                    id,
		logger:                logger,
		staticVideosDirectory: staticVideosDirectory,
		thumbnailsPath:        thumbnailsPath,
		streamURL:             streamURL,
	}
}

func (vr *VideoRecorder) StartRecording() error {
	vr.mu.Lock()
	defer vr.mu.Unlock()
	fmt.Printf("start recording with cameraURL: %s\n", vr.streamURL)

	if vr.isRecording {
		log.Println(fmt.Errorf("recording is already in progress"))
		return fmt.Errorf("recording is already in progress")
	}

	dir, err := utils.CreateDataYMD(vr.staticVideosDirectory, vr.ID)
	if err != nil {
		vr.logger.Println(err)
		return fmt.Errorf("error with create YMD-data: %w", err)
	}

	outputFile := "%H:%M:%S"
	vr.lastVideoFile = fmt.Sprintf("%s.%s", time.Now().Format("15:04:05.00"), "mp4")

	var stderr bytes.Buffer
	vr.cmd = exec.Command("ffmpeg",
		"-i", vr.streamURL, // URL камеры
		"-c:v", "libx264", // Копируем видеопоток без перекодирования
		"-b:v", "1500k",
		"-an",           // Копируем аудиопоток без перекодирования
		"-f", "segment", // Формат сегментации
		"-segment_time", "100", // Длительность сегмента — 60 секунд
		"-reset_timestamps", "1", // Сброс таймстампов для каждого сегмента
		"-strftime", "1",
		fmt.Sprintf("%s/%s.mp4", dir, outputFile),
	)
	vr.cmd.Stderr = &stderr

	if err := vr.cmd.Start(); err != nil {
		fmt.Printf("Вывод FFmpeg: %s", stderr.String())
		vr.logger.Println(stderr.String())
		return fmt.Errorf("failed to start recording: %w", err)
	}

	go vr.generateThumbnailFromStream(vr.streamURL, vr.lastVideoFile)
	go vr.monitorRecording(&stderr)

	vr.isRecording = true

	return nil
}

func (vr *VideoRecorder) monitorRecording(stderr *bytes.Buffer) {
	if err := vr.cmd.Wait(); err != nil {
		vr.logger.Printf("Error during recording: %v, FFmpeg output: %s", err, stderr.String())
	}
	vr.mu.Lock()
	vr.isRecording = false
	vr.mu.Unlock()

}

// StopRecording останавливает процесс записи видео
func (vr *VideoRecorder) StopRecording() error {
	vr.mu.Lock()
	defer vr.mu.Unlock()

	if !vr.isRecording {
		return fmt.Errorf("no recording in progress")
	}

	if err := vr.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("error stopping recording: %w", err)
	}

	//wait не выполняется, поэтому освобождаем ресурсы с помощьюRelease
	if err := vr.cmd.Process.Release(); err != nil {
		vr.logger.Printf("Failed to RELEASE FFmpeg proccess: %s", err)
		return fmt.Errorf("error waiting for ffmpeg to finish: %w", err)
	}

	vr.isRecording = false
	return nil
}

func (vr *VideoRecorder) generateThumbnailFromStream(cameraURL, filename string) {
	dir, err := utils.CreateDataYMD(vr.thumbnailsPath, vr.ID)
	if err != nil {
		vr.logger.Printf("error witch create ")
	}
	vr.lastVideoFile = fmt.Sprintf("%s.%s", time.Now().Format("15:04:05.00"), "mp4")
	cmd := exec.Command("ffmpeg",
		"-i", cameraURL,
		"-ss", "00:00:00.500",
		"-vframes", "1",
		fmt.Sprintf("%s/%s.png", dir, filename),
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("err with creating thumbnails: %s\n", err)
	}

}
