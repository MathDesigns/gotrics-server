package api

import (
	"encoding/json"
	"gotrics-server/internal/storage"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) handlePostMetrics(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	if tokenParts[1] != s.config.AgentAuthToken {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
		return
	}

	var metric storage.Metric
	if err := c.ShouldBindJSON(&metric); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.store.WriteMetric(c.Request.Context(), &metric); err != nil {
		s.logger.Printf("Failed to write metric to store: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store metric"})
		return
	}

	metricJSON, err := json.Marshal(metric)
	if err == nil {
		s.hub.broadcast <- metricJSON
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "ok"})
}

func (s *Server) handleGetMetrics(c *gin.Context) {
	hostname := c.Param("hostname")
	durationStr := c.DefaultQuery("last", "1h")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid duration format"})
		return
	}

	metrics, err := s.store.GetHostMetrics(c.Request.Context(), hostname, duration)
	if err != nil {
		s.logger.Printf("Failed to get host metrics from store: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve metrics"})
		return
	}

	c.JSON(http.StatusOK, metrics)
}
