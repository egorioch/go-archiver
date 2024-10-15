package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-archiver/internal/config_"
	"go-archiver/internal/logger"
	"go-archiver/internal/recording/service"
	"log"
)

type MultipleRecordingHandler struct {
	URLs                  []string
	StreamRecordingLogger *log.Logger
	RecordedVideoPath     string
	ThumbnailsPath        string
	MultipleRecordingUnit *service.MultipleRecording
}

func InitMultipleRecordingHandler() *MultipleRecordingHandler {
	config, err := config_.ConfigLoader("static/config.yaml")
	if err != nil {
		fmt.Printf("cannot open cfg-file: %s", err)
	}

	streamLogger, err := logger.CreateLogger(config.Logs.Videos)
	recordedVideoPath := config.StaticPaths.VideosPath
	thumbnailsPath := config.StaticPaths.ThumbnailsPath
	urls := config.Stream.URLs

	// указатель на один экзмепляр для Start и для Stop(recording)
	mru := service.InitMultipleRecord(urls, streamLogger, recordedVideoPath, thumbnailsPath)

	return &MultipleRecordingHandler{
		URLs:                  urls,
		StreamRecordingLogger: streamLogger,
		RecordedVideoPath:     recordedVideoPath,
		ThumbnailsPath:        thumbnailsPath,
		MultipleRecordingUnit: mru,
	}
}

func (mrh *MultipleRecordingHandler) StartMultipleRecording(c *gin.Context) {
	err := mrh.MultipleRecordingUnit.StartMultipleCameras()
	if err != nil {
		mrh.StreamRecordingLogger.Printf("error with start record with all cameras: %v", err)
		c.JSON(500, gin.H{"error": err})
	}

	c.JSON(200, gin.H{"message": "Recording started"})
}

// наверное нужно создать для start и stop один экземпляр MRH в app.go
func (mrh *MultipleRecordingHandler) StopMultipleRecording(c *gin.Context) {
	//mr := service.InitMultipleRecord(mrh.URLs, mrh.StreamRecordingLogger, mrh.RecordedVideoPath, mrh.ThumbnailsPath)
	err := mrh.MultipleRecordingUnit.StopMultipleCameras()
	if err != nil {
		mrh.StreamRecordingLogger.Printf("error with STOP record with all cameras: %v", err)
		c.JSON(500, gin.H{"error": err})
	}

	c.JSON(200, gin.H{"message": "Recording stopped"})
}

func CreateMultiRecHandler(r *gin.Engine) {
	mrh := InitMultipleRecordingHandler()

	r.POST("/api/start_multiply_recording", mrh.StartMultipleRecording)
	r.POST("/api/stop_multiply_recording", mrh.StopMultipleRecording)

}
