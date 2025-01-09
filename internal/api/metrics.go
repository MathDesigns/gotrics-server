package api

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type MetricsStore struct {
	sync.RWMutex
	Data map[string]Metrics
}

var store = MetricsStore{
	Data: make(map[string]Metrics),
}

type Metrics struct {
	NodeID      string `json:"node_id"`
	CPU         string `json:"cpu"`
	CPUModel    string `json:"cpu_model"`
	NumCores    int    `json:"num_cores"`
	NumThreads  int    `json:"num_threads"`
	UsedMemory  uint64 `json:"used_memory"`
	TotalMemory uint64 `json:"total_memory"`
	Uptime      uint64 `json:"uptime"`
	Platform    string `json:"platform"`
	UsedSpace   uint64 `json:"used_space"`
	TotalSpace  uint64 `json:"total_space"`
}

func handlePostMetrics(c *gin.Context) {
	var metrics Metrics
	if err := c.ShouldBindJSON(&metrics); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	store.Lock()
	store.Data[metrics.NodeID] = metrics
	store.Unlock()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func handleGetMetrics(c *gin.Context) {
	store.RLock()
	defer store.RUnlock()

	metrics := make([]Metrics, 0, len(store.Data))
	for _, m := range store.Data {
		metrics = append(metrics, m)
	}
	c.JSON(http.StatusOK, metrics)
}

func RegisterRoutes(router *gin.Engine) {
	router.POST("/metrics", handlePostMetrics)
	router.GET("/metrics", handleGetMetrics)
}
