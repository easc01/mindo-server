package communityhandler

import (
	"net/http"
	"sync"

	"github.com/easc01/mindo-server/internal/middleware"
	"github.com/easc01/mindo-server/internal/models"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	httputil "github.com/easc01/mindo-server/pkg/utils/http_util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type RoomManager struct {
	rooms map[string]map[*dto.ChatClient]bool
	mu    sync.Mutex
}

// Map<RoomID, ChatClientDTO>
var roomManager = RoomManager{
	rooms: make(map[string]map[*dto.ChatClient]bool),
}

// AddClient registers a new client in a specific room
func (m *RoomManager) AddClient(client *dto.ChatClient, roomID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// If the room doesn't exist, create it
	if _, exists := m.rooms[roomID]; !exists {
		m.rooms[roomID] = make(map[*dto.ChatClient]bool)
	}

	// Add the client to the room
	m.rooms[roomID][client] = true
}

// RemoveClient removes a client from a specific room
func (m *RoomManager) RemoveClient(client *dto.ChatClient, roomID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Remove the client from the room
	if roomClients, exists := m.rooms[roomID]; exists {
		delete(roomClients, client)
		if len(roomClients) == 0 {
			// If no clients left in the room, delete the room
			delete(m.rooms, roomID)
		}
	}

	client.Conn.Close()
}

func HandleRoomChatWS(c *gin.Context) {

	// Http part, do auth and stuff
	r := c.Request
	w := c.Writer.(http.ResponseWriter)

	roomID := r.URL.Query().Get("communityId")
	authToken := r.URL.Query().Get("auth")

	r.Header.Set(constant.Authorization, authToken)
	user, err := middleware.AuthenticateAndFetchUser(r, models.UserTypeAppUser)

	if err != nil {
		httputil.NewErrorResponse(http.StatusUnauthorized, err.Error(), nil).Send(c)
		c.Abort()
		return
	}

	// TODO: check whether user is in community

	// Once Authenticated, Only then Upgrade to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Log.Errorf("upgrade error, %s", err)
		return
	}
	logger.Log.Infof("client connected to room, %s", roomID)

	client := dto.ChatClient{
		AppUser: user.AppUser,
		Conn:    conn,
	}

	// Add client to the room
	roomManager.AddClient(&client, roomID)
	defer roomManager.RemoveClient(&client, roomID)

	// Handle incoming messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Log.Errorf("read error, %s", err)
			break
		}
		logger.Log.Infof("received in room %s: %s", roomID, msg)
		// Broadcast message to the room
		roomManager.Broadcast(roomID, msg)
	}
}

// Broadcast sends a message to all clients in a specific room
func (m *RoomManager) Broadcast(roomID string, msg []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the room exists
	if roomClients, exists := m.rooms[roomID]; exists {
		for client := range roomClients {
			err := client.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				logger.Log.Errorf("write error, %s", err)
				client.Conn.Close()
				delete(roomClients, client)
			}
		}
	}
}
