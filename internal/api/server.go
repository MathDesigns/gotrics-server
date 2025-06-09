package api

import (
	"gotrics-server/internal/config"
	"gotrics-server/internal/storage"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	logger *log.Logger
	config *config.Config
	store  *storage.InfluxStore
	hub    *Hub
	router *gin.Engine
}

func NewServer(logger *log.Logger, cfg *config.Config, store *storage.InfluxStore, hub *Hub) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	server := &Server{
		logger: logger,
		config: cfg,
		store:  store,
		hub:    hub,
		router: router,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	apiV1 := s.router.Group("/api/v1")
	{
		apiV1.POST("/metrics", s.handlePostMetrics)
		apiV1.GET("/metrics/:hostname", s.handleGetMetrics)
		apiV1.GET("/nodes", s.handleGetNodes)

		apiV1.POST("/hardware", s.handlePostHardware)
		apiV1.GET("/hardware/:hostname", s.handleGetHardware)
	}
	s.router.GET("/ws", s.handleWebSocket)
}

func (s *Server) Router() *gin.Engine {
	return s.router
}
