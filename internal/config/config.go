package config

import (
	"os"
	"strconv"
	"time"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	JWT      JWTConfig      `json:"jwt"`
	Admin    AdminConfig    `json:"admin"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver string `json:"driver"`
	DSN    string `json:"dsn"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string        `json:"secret"`
	Expiration time.Duration `json:"expiration"`
}

// AdminConfig 管理员配置
type AdminConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Load 加载配置
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getPortEnv(),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			Driver: getEnv("DB_DRIVER", "postgres"),
			DSN:    getEnv("DB_DSN", "host=localhost port=5432 user=goblog password=goblog123 dbname=goblog sslmode=disable"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key"),
			Expiration: getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
		},
		Admin: AdminConfig{
			Username: getEnv("ADMIN_USERNAME", "admin"),
			Password: getEnv("ADMIN_PASSWORD", "admin123"),
		},
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getPortEnv 获取端口环境变量，支持PORT和SERVER_PORT
func getPortEnv() string {
	// 优先检查标准的PORT环境变量
	if port := os.Getenv("PORT"); port != "" {
		// 如果PORT不以冒号开头，添加冒号
		if port[0] != ':' {
			return ":" + port
		}
		return port
	}

	// 检查自定义的SERVER_PORT环境变量
	if port := os.Getenv("SERVER_PORT"); port != "" {
		return port
	}

	// 默认端口
	return ":8080"
}

// getDurationEnv 获取时间间隔环境变量
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// getIntEnv 获取整数环境变量
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
