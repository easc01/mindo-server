package communityhandler

import (
	"context"
	"net/http"
	"sync"

	"github.com/easc01/mindo-server/internal/middleware"
	"github.com/easc01/mindo-server/internal/models"
	communityservice "github.com/easc01/mindo-server/internal/services/community_service"
	"github.com/easc01/mindo-server/pkg/dto"
	"github.com/easc01/mindo-server/pkg/logger"
	"github.com/easc01/mindo-server/pkg/utils/constant"
	networkutil "github.com/easc01/mindo-server/pkg/utils/network_util"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type RoomManager struct {
	rooms map[uuid.UUID]map[*dto.ChatClient]bool
	mu    sync.Mutex
}

// Map<RoomID, ChatClientDTO>
var roomManager = RoomManager{
	rooms: make(map[uuid.UUID]map[*dto.ChatClient]bool),
}

// AddClient registers a new client in a specific room
func (m *RoomManager) AddClient(client *dto.ChatClient, roomID uuid.UUID) {
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
func (m *RoomManager) RemoveClient(client *dto.ChatClient, roomID uuid.UUID) {
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

func HandleRoomChatWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Log.Errorf("upgrade error, %s", err)
		networkutil.WSError(
			http.StatusInternalServerError,
			err.Error(),
		).Send(conn)
		return
	}

	roomID := r.URL.Query().Get("communityId")
	authToken := r.URL.Query().Get("auth")

	parsedRoomID, err := uuid.Parse(roomID)
	if err != nil {
		networkutil.WSError(
			http.StatusBadRequest,
			"invalid community id",
		).Send(conn)
		return
	}

	// authenticate connection
	r.Header.Set(constant.Authorization, authToken)
	user, err := middleware.AuthenticateAndFetchUser(r, models.UserTypeAppUser)

	if err != nil {
		networkutil.WSError(
			http.StatusUnauthorized,
			err.Error(),
		).Send(conn)
		return
	}

	// check whether user is in community
	var userJoinedCommunity *dto.CommunityDTO
	for _, community := range user.AppUser.JoinedCommunities {
		if community.ID == parsedRoomID {
			userJoinedCommunity = &community
		}
	}

	if userJoinedCommunity == nil {
		networkutil.WSError(
			http.StatusNotFound,
			"community not found",
		).Send(conn)
		return
	}

	client := dto.ChatClient{
		AppUser: user.AppUser,
		Conn:    conn,
	}

	// Add client to the room
	roomManager.AddClient(&client, parsedRoomID)
	defer roomManager.RemoveClient(&client, parsedRoomID)

	// Handle incoming messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logger.Log.Errorf("read error, %s", err)
			break
		}
		logger.Log.Infof("received in room %s: %s", parsedRoomID, msg)
		// Broadcast message to the room
		roomManager.Broadcast(parsedRoomID, user.AppUser.UserID, msg)
	}
}

// Broadcast sends a message to all clients in a specific room
func (m *RoomManager) Broadcast(roomID uuid.UUID, userID uuid.UUID, msg []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the room exists
	if roomClients, exists := m.rooms[roomID]; exists {
		_, err := communityservice.SaveCommunityMessage(
			context.Background(),
			roomID,
			userID,
			string(msg),
		)

		if err != nil {
			logger.Log.Errorf("failed to save message, %s", err.Error())
			return
		}

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
