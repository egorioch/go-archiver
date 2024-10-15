package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"go-archiver/internal/recording/handlers"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()
	customPrometheusRegistry := prometheus.NewRegistry()

	//handlers.CreateVideoRecordingHandler(r, customPrometheusRegistry)
	handlers.CreateMultiRecHandler(r)
	handlers.CreateVideoSenderHandler(r, customPrometheusRegistry)
	r.GET("/metrics", prometheusHandler(customPrometheusRegistry))

	c := createCors()
	httpHandler := c.Handler(r)
	// Запуск сервера
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", httpHandler); err != nil {
		log.Fatal(err)
	}
}

func prometheusHandler(rgstr *prometheus.Registry) gin.HandlerFunc {
	h := promhttp.HandlerFor(rgstr, promhttp.HandlerOpts{})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func createCors() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Разрешенные источники
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
	})
}
