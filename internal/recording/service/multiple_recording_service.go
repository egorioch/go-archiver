package service

import (
	"fmt"
	"log"
	"sync"
)

// Сервис, позволяющий запускать запись m-rtsp потоков
type MultipleRecording struct {
	Streams               *[]Stream
	StreamRecordingLogger *log.Logger
	RecorderStreamPath    string
	ThumbnailsPath        string
}

func InitMultipleRecord(URLs []string, streamRecordingLogger *log.Logger, recordedStreamPath, thumbnailsPath string) *MultipleRecording {
	var streams []Stream
	for i, s := range URLs {
		vr := NewVideoRecorder(i, streamRecordingLogger, recordedStreamPath, thumbnailsPath, s)
		streams = append(streams, *InitNewStream(i, vr))
	}
	return &MultipleRecording{
		Streams:               &streams,
		StreamRecordingLogger: streamRecordingLogger,
		ThumbnailsPath:        thumbnailsPath,
	}
}

// StartMultipleCameras Запуск нескольких камер.
func (mr *MultipleRecording) StartMultipleCameras() error {
	var wg sync.WaitGroup

	for _, stream := range *mr.Streams {
		wg.Add(1)
		streamCopy := stream
		go func(stream_ *Stream) {
			defer wg.Done()

			if err := stream_.VideoRecorder.StartRecording(); err != nil {
				//stream_.VideoRecorder.isRecording = true

				mr.StreamRecordingLogger.Printf("Failed to start recording for camera[%d] %s: %v", stream_.ID, stream_.VideoRecorder.streamURL, err)
				fmt.Printf("err with stream[%d]): %s \n", stream_.ID, err)
			}
		}(&streamCopy)
	}

	wg.Wait()
	return nil
}

func (mr *MultipleRecording) StopMultipleCameras() error {
	var wg sync.WaitGroup

	for _, stream := range *mr.Streams {
		if stream.VideoRecorder.isRecording {
			wg.Add(1)
			go func(vr VideoRecorder) {
				defer wg.Done()

				if err := vr.StopRecording(); err != nil {
					mr.StreamRecordingLogger.Printf("Failed to stop recording for camera %s: %v", vr.streamURL, err)
					fmt.Printf("err with stream[%s]): %v \n", vr.streamURL, err)
				}
			}(*stream.VideoRecorder)
		}
	}

	wg.Wait()
	return nil
}
