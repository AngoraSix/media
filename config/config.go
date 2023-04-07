package config

import (
	"os"
	"strconv"
)

// GetEnv gets an environment variable and provides a fallback
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Simple helper function to read an environment variable into unsigned integer or return a default value
func getEnvAsUint(name string, defaultVal uint) uint {
	valueStr := GetEnv(name, "")
	if value, err := strconv.ParseUint(valueStr, 10, 32); err == nil {
		return uint(value)
	}
	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := GetEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// ServerConfig ...
type ServerConfig struct {
	Hostname string
	Port     string
	Scheme   string
}

// ServiceConfig ...
type ServiceConfig struct {
	Strategy             string
	StaticDir            string
	UploadDir            string
	BucketName           string
	ProjectID            string
	StorageAPIHost       string
	ThumbnailMaxWidth    uint
	ThumbnailMaxHeight   uint
	ThumbnailMaxSize     uint
	ThumbnailQuality     uint
	ThumbnailCompression int
}

// GetHost ...
func (sc *ServerConfig) GetHost() string {
	return sc.Scheme + "://" + sc.Hostname + ":" + sc.Port
}

// DefaultServerConfig instance
var DefaultServerConfig = ServerConfig{
	GetEnv("A6_MEDIA_SVC_HOSTNAME", "0.0.0.0"),
	GetEnv("A6_MEDIA_SVC_PORT", "80"),
	GetEnv("A6_MEDIA_SVC_PROTOCOL", "http"),
}

// DefaultServiceConfig instance
var DefaultServiceConfig = ServiceConfig{
	GetEnv("A6_MEDIA_SVC_STRATEGY", "local"),
	GetEnv("A6_MEDIA_SVC_STATIC_DIR", "static"),
	GetEnv("A6_MEDIA_SVC_UPLOADS_DIR", "uploads"),
	GetEnv("A6_MEDIA_SVC_BUCKET_NAME", "hoc-storage"),
	GetEnv("A6_MEDIA_SVC_PROJECT_ID", "angorasix-203314"),
	GetEnv("A6_MEDIA_SVC_STORAGE_API_HOST", "https://storage.googleapis.com"),
	getEnvAsUint("A6_MEDIA_SVC_THUMBNAIL_MAX_WIDTH", 800),
	getEnvAsUint("A6_MEDIA_SVC_THUMBNAIL_MAX_HEIGHT", 800),
	// max size in bytes (B)
	getEnvAsUint("A6_MEDIA_SVC_THUMBNAIL_MAX_SIZE", 600000),
	// quality to create the thumbnail (0-100)
	getEnvAsUint("A6_MEDIA_SVC_THUMBNAIL_JPG_QUALITY", 90),
	// compression level:
	// DefaultCompression (0)
	// NoCompression (-1)
	// BestSpeed (-2)
	// BestCompression (-3)
	getEnvAsInt("A6_MEDIA_SVC_THUMBNAIL_PNG_COMPRESSION", 0),
}
