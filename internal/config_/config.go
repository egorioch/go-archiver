package config_

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// Структура для хранения данных из YAML
type Config struct {
	StaticPaths struct {
		VideosPath     string `yaml:"videos_path"`
		ThumbnailsPath string `yaml:"thumbnails_path"`
	} `yaml:"static_paths"`
	Logs struct {
		Videos     string `yaml:"videos"`
		Thumbnails string `yaml:"thumbnails"`
		FFmpeg     string `yaml:"ffmpeg"`
	} `yaml:"logs"`
	Stream struct {
		URLs []string `yaml:"urls"`
	} `yaml:"stream"`

	VideoLogger  *log.Logger
	FFmpegLogger *log.Logger
}

// ConfigLoader возвращает экземпляр, с загружеными yaml-данными
func ConfigLoader(configPath string) (*Config, error) {
	config := &Config{}

	err := validateConfigPath(configPath)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error with open cfg-path: %w", err)
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func createLogger(filePath string) (*log.Logger, *os.File, error) {
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	logger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	return logger, logFile, nil
}
