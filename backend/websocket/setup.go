package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	// "os"
	"sync"
	"time"
	"wechat/models"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

func InitWebSocketDB(database *gorm.DB) {
	db = database
}

var (
	roomClients   = make(map[uuid.UUID]map[*websocket.Conn]bool)
	roomClientsMu sync.Mutex
	db            *gorm.DB
)

type WSMessage struct {
	Type string `json:"type"` // "join" | "message"
	Room string `json:"room"`
	User string `json:"user"`
	Content string `json:"content"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// origin := r.Header.Get("Origin")
		// allowed := []string{
		// 	os.Getenv("FRONTEND_ORIGIN"),
		// 	"http://localhost:3000",
		// 	"https://wechat-livid.vercel.app/",
		// 	"https://wechat-production-2d6e.up.railway.app/",
		// 	"http://127.0.0.1:3000",
		// 	"http://192.168.1.157:3000",
		// }
		// for _, o := range allowed {
		// 	if origin == o && origin != "" {
		// 		return true
		// 	}
		// }
		// log.Printf("[WebSocket] Rejected origin: %s", origin)
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	var currentRoomID uuid.UUID
	var currentUserID uuid.UUID
	conn.SetPongHandler(func(appData string) error {
		return nil
	})
	go keepAlive(conn)

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Println("Invalid message JSON:", err)
			continue
		}

		roomID, err := uuid.Parse(msg.Room)
		userID, err2 := uuid.Parse(msg.User)

		if err != nil || err2 != nil {
			log.Println("Invalid UUIDs in join/message")
			continue
		}
		switch msg.Type {
		case "join":
			currentRoomID = roomID
			currentUserID = userID

			addClientToRoom(roomID, conn)

			// Preload messages with User info
			var messages []models.Message
			db.Preload("User").
				Where("room_id = ?", roomID).
				Order("created_at asc").
				Find(&messages)

			// Send as a single "history" event
			history := struct {
				Type     string           `json:"type"`
				Messages []models.Message `json:"messages"`
			}{
				Type:     "history",
				Messages: messages,
			}

			jsonHistory, _ := json.Marshal(history)
			conn.WriteMessage(websocket.TextMessage, jsonHistory)

			log.Printf("User %s joined room %s\n", userID, roomID)
		case "message":
			if currentRoomID == uuid.Nil || currentUserID == uuid.Nil {
				log.Printf("Message sent before joining room")
				continue
			}
			chat := models.Message{
				UserID:  &userID,
				RoomID:  &roomID,
				Content: msg.Content,
			}
			if err := db.Create(&chat).Error; err != nil {
				log.Panicln("DB insert failed: ", err)
				continue
			}
			db.Preload("User").First(&chat, "id = ?", chat.ID)

			broadcastToRoom(roomID, chat)

		}

	}
	if currentRoomID != uuid.Nil {
		removeClientFromRoom(currentRoomID, conn)
		log.Printf("User %s left the room %s /n", currentUserID, currentRoomID)
	}
}
func keepAlive(conn *websocket.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Println("[WebSocket]  Ping failed, closing connection:", err)
			conn.Close()
			return
		}
	}
}
func addClientToRoom(roomID uuid.UUID, conn *websocket.Conn) {
	roomClientsMu.Lock()
	defer roomClientsMu.Unlock()
	if roomClients[roomID] == nil {
		roomClients[roomID] = make(map[*websocket.Conn]bool)
	}
	roomClients[roomID][conn] = true
}
func removeClientFromRoom(roomID uuid.UUID, conn *websocket.Conn) {
	roomClientsMu.Lock()
	defer roomClientsMu.Unlock()
	if clients, ok := roomClients[roomID]; ok {
		delete(clients, conn)
		conn.Close()
		if len(clients) == 0 {
			delete(roomClients, roomID)
		}
	}

}
func broadcastToRoom(roomID uuid.UUID, msg models.Message) {
	roomClientsMu.Lock()
	clients := make([]*websocket.Conn, 0, len(roomClients[roomID]))
	for conn := range roomClients[roomID] {
		clients = append(clients, conn)
	}
	roomClientsMu.Unlock()

	wrapped := struct {
		Type    string         `json:"type"`
		Message models.Message `json:"message"`
	}{
		Type:    "message",
		Message: msg,
	}

	data, _ := json.Marshal(wrapped)
	for _, conn := range clients {
		go func(c *websocket.Conn) {
			if err := c.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("Broadcast failed: ", err)
				removeClientFromRoom(roomID, c)
			}
		}(conn)
	}
}
