package api

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// MetricsStore to hold the latest metrics data
type MetricsStore struct {
	sync.RWMutex
	Data map[string]Metrics
}

var store = MetricsStore{
	Data: make(map[string]Metrics),
}

// Metrics struct
type Metrics struct {
	NodeID string  `json:"node_id"`
	CPU    float64 `json:"cpu"`
	Memory uint64  `json:"memory"`
}

// handlePostMetrics handles the POST request from nodes
func handlePostMetrics(c *gin.Context) {
	var metrics Metrics
	if err := c.ShouldBindJSON(&metrics); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save metrics in store
	store.Lock()
	store.Data[metrics.NodeID] = metrics
	store.Unlock()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// handleGetMetrics handles the GET request from the front-end
func handleGetMetrics(c *gin.Context) {
	store.RLock()
	defer store.RUnlock()

	// Return all metrics as a list
	metrics := make([]Metrics, 0, len(store.Data))
	for _, m := range store.Data {
		metrics = append(metrics, m)
	}
	c.JSON(http.StatusOK, metrics)
}

// RegisterRoutes registers the routes for the server
func RegisterRoutes(router *gin.Engine) {
	router.POST("/metrics", handlePostMetrics)
	router.GET("/metrics", handleGetMetrics)
}
