package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go-archiver/internal/config_"
	"go-archiver/internal/logger"
	"go-archiver/internal/recording/service"
	"go-archiver/package/custom_prometheus"
	"log"
)

type VideoRecordingHandler struct {
	CameraURL                 string
	VideoRecorder             *service.VideoRecorder
	Logger                    *log.Logger
	VideoRecorderErrorCounter *prometheus.CounterVec
}

func (vrh *VideoRecordingHandler) httpRecordError(statusCode, errorType string) {
	vrh.VideoRecorderErrorCounter.With(prometheus.Labels{"status_code": statusCode, "error_type": errorType}).Inc()
}

func initVideoHandler(cameraURL string, registry *prometheus.Registry) *VideoRecordingHandler {
	config, err := config_.ConfigLoader("static/config.yaml")
	cc := custom_prometheus.CreateCustomCounter(
		"video_record_handler",
		"Total number of errors encountered by VideoSender, categorized by error type",
		[]string{"status_code", "error_type"},
	)
	registry.MustRegister(cc)

	if err != nil {
		fmt.Printf("cannot open cfg-file: %s", err)
	}

	streamLogger, err := logger.CreateLogger(config.Logs.Videos)
	//ffmpegLogger, err := logger.CreateLogger(config.Logs.FFmpeg)
	//videosPath := config.StaticPaths.VideosPath
	//thumbsPath := config.StaticPaths.ThumbnailsPath

	return &VideoRecordingHandler{
		CameraURL:                 cameraURL,
		VideoRecorder:             nil,
		Logger:                    streamLogger,
		VideoRecorderErrorCounter: cc,
	}
}

func (vrh *VideoRecordingHandler) startRecording(c *gin.Context) {
	err := vrh.VideoRecorder.StartRecording()
	if err != nil {
		vrh.Logger.Println(err)
		vrh.httpRecordError("400", err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Recording started"})
}

func (vrh *VideoRecordingHandler) stopRecording(c *gin.Context) {
	err := vrh.VideoRecorder.StopRecording()
	if err != nil {
		vrh.Logger.Println(err)
		vrh.httpRecordError("400", err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Recording completed"})
}

func CreateVideoRecordingHandler(r *gin.Engine, registry *prometheus.Registry, cameraURL string) {
	ivh := initVideoHandler(cameraURL, registry)

	r.POST("/api/start_recording", ivh.startRecording)
	r.POST("/api/stop_recording", ivh.stopRecording)
}
