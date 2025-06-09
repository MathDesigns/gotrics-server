package api

import (
	"encoding/json"
	"gotrics-server/internal/hardware"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const hardwareDataDir = "data/details"

func (s *Server) handlePostHardware(c *gin.Context) {
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

	var info hardware.HardwareInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := os.MkdirAll(hardwareDataDir, 0755); err != nil {
		s.logger.Printf("Failed to create hardware data directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create data directory"})
		return
	}

	filePath := filepath.Join(hardwareDataDir, info.Hostname+".json")
	file, err := os.Create(filePath)
	if err != nil {
		s.logger.Printf("Failed to create hardware info file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create data file"})
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(info); err != nil {
		s.logger.Printf("Failed to write hardware info to file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) handleGetHardware(c *gin.Context) {
	hostname := c.Param("hostname")
	filePath := filepath.Join(hardwareDataDir, hostname+".json")
	c.File(filePath)
}
