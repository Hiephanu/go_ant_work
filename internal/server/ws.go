package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebsocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) UpgradeToWS(c *gin.Context) {
	roomID := c.Query("roomId")
	token := c.Query("token")

	if roomID == "" || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing token or room id"})
		return
	}

	tokenData, err := DecryptToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading to websocket:", err)
		return
	}
	defer conn.Close()

	if err := s.roomTracker.AddConnToUser(tokenData.UserID, roomID, conn); err != nil {
		fmt.Printf("Error adding connection to user %s in room %s: %v\n", tokenData.UserID, roomID, err)
		conn.Close()
		return
	}

	fmt.Println("WebSocket connected for user:", tokenData.UserID)

	go s.handleWebSocket(conn, tokenData.UserID, roomID)
}

func (s *Server) handleWebSocket(conn *websocket.Conn, userID string, roomID string) {
	defer conn.Close()

	for {
		fmt.Println(conn)
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		var wsMessage WebsocketMessage
		if err = json.Unmarshal(message, &wsMessage); err != nil {
			fmt.Println("Error unmarshalling message:", err)
			break
		}

		fmt.Printf("Received message from user %s in room %s: %s\n", userID, roomID, wsMessage.Type)

		if err = conn.WriteMessage(messageType, message); err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}
