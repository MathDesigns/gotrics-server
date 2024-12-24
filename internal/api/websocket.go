package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gotrics-server/internal/logger"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			logger.Error.Printf("Failed to close websocket connection: %v", err)
		}
	}(conn)

	clients[conn] = true

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				err := client.Close()
				if err != nil {
					logger.Error.Printf("Failed to close websocket connection: %v", err)
				}
				delete(clients, client)
			}
		}
	}
}
