package health

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

type HealthResponse struct {
	Status      string     `json:"status"`
	Timestamp   int64      `json:"timestamp"`
	Environment string     `json:"environment"`
	System      SystemInfo `json:"system"`
}

type SystemInfo struct {
	Version      string `json:"version"`
	GoVersion    string `json:"go_version"`
	NumGoroutine int    `json:"goroutines"`
	MemoryUsage  string `json:"memory_usage"`
	NumCPU       int    `json:"num_cpu"`
	Uptime       string `json:"uptime"`
}

var startTime = time.Now()

func (h *Handler) Health(c *gin.Context) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	systemInfo := SystemInfo{
		Version:      "1.0.0", // You can make this configurable
		GoVersion:    runtime.Version(),
		NumGoroutine: runtime.NumGoroutine(),
		MemoryUsage:  formatBytes(mem.Alloc),
		NumCPU:       runtime.NumCPU(),
		Uptime:       formatUptime(time.Since(startTime)),
	}

	response := HealthResponse{
		Status:      "OK",
		Timestamp:   time.Now().Unix(),
		Environment: getEnvironment(),
		System:      systemInfo,
	}

	c.JSON(http.StatusOK, response)
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatUptime(d time.Duration) string {
	d = d.Round(time.Second)
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	seconds := d / time.Second

	if days > 0 {
		return fmt.Sprintf("%dd%dh%dm%ds", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh%dm%ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

func getEnvironment() string {
	env := gin.Mode()
	if env == gin.ReleaseMode {
		return "production"
	}
	if env == gin.TestMode {
		return "test"
	}
	return "development"
}
