package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) handleGetNodes(c *gin.Context) {
	hosts, err := s.store.GetKnownHosts(c.Request.Context())
	if err != nil {
		s.logger.Printf("Failed to get known hosts from store: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve hosts"})
		return
	}

	c.JSON(http.StatusOK, hosts)
}
