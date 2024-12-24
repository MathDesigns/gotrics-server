package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gotrics-server/internal/config"
	"gotrics-server/internal/logger"
)

func StartServer(cfg *config.Config) {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	RegisterRoutes(r)

	serverAddr := ":" + cfg.ServerPort
	logger.Info.Printf("Starting server on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		logger.Info.Fatalf("Error starting server: %v", err)
	}
}
