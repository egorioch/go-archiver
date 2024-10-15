package handlers

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go-archiver/package/custom_prometheus"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type VideoSender struct {
	VideoSenderErrorCounter *prometheus.CounterVec
}

func initVideoSender(registry *prometheus.Registry) *VideoSender {
	cc := custom_prometheus.CreateCustomCounter(
		"video_sender_errors_total",
		"Total number of errors encountered by VideoSender, categorized by error type",
		[]string{"status_code", "error_type"},
	)
	registry.MustRegister(cc)

	return &VideoSender{
		VideoSenderErrorCounter: cc,
	}
}

func (vs *VideoSender) videoSenderPrometheusError(statusCode, errorType string) {
	vs.VideoSenderErrorCounter.With(prometheus.Labels{"status_code": statusCode, "error_type": errorType}).Inc()
}

func (vs *VideoSender) listVideos(c *gin.Context) {
	//videoDir := "./static/videos"
	thumbnailDir := "./static/thumbnails"

	var videos []map[string]string

	err := filepath.Walk(thumbnailDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Проверяем, является ли текущий файл файлом с расширением .png
		if !info.IsDir() && filepath.Ext(info.Name()) == ".png" {
			fmt.Println("Найден файл:", path)
			thumbnailData, err := ioutil.ReadFile(path)
			var thumbnailBase64 string
			if err == nil {
				// Конвертируем превью в base64
				thumbnailBase64 = base64.StdEncoding.EncodeToString(thumbnailData)
			} else {
				thumbnailBase64 = ""
			}

			filenameOnly := filepath.Base(path)
			videos = append(videos, map[string]string{
				"name":      filenameOnly,
				"thumbnail": thumbnailBase64,
			})
		}

		return nil
	})

	if err != nil {
		fmt.Println("Ошибка обхода директорий:", err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not list video files"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"videos": videos})
}
func (vs *VideoSender) sendVideoFile(c *gin.Context) {
	// Предположим, что видео хранится в директории ./videos
	videoFile := c.Param("filename")
	videoPath := filepath.Join("./static/videos", videoFile)
	fmt.Println("requested filename: ", videoFile)

	// Проверка наличия файла
	if _, err := filepath.Abs(videoPath); err != nil {
		vs.videoSenderPrometheusError("404", err.Error())
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}

	// Установка заголовка Content-Type для видео
	c.Header("Content-Type", "videos/mp4")

	// Поддержка Range-запросов
	c.File(videoPath)
}

func CreateVideoSenderHandler(r *gin.Engine, registry *prometheus.Registry) {
	vs := initVideoSender(registry)

	r.GET("/api/list_videos", vs.listVideos)
	r.POST("/api/:filename", vs.sendVideoFile)
}
