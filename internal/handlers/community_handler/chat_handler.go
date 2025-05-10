package communityhandler

import (
	"context"
	"encoding/json"
	"fmt"
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

var roomManager = RoomManager{
	rooms: make(map[uuid.UUID]map[*dto.ChatClient]bool),
}

func (m *RoomManager) AddClient(client *dto.ChatClient, roomID uuid.UUID) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.rooms[roomID]; !exists {
		m.rooms[roomID] = make(map[*dto.ChatClient]bool)
	}
	m.rooms[roomID][client] = true
}

func (m *RoomManager) RemoveClient(client *dto.ChatClient, roomID uuid.UUID) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if roomClients, exists := m.rooms[roomID]; exists {
		delete(roomClients, client)
		if len(roomClients) == 0 {
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

	r.Header.Set(constant.Authorization, authToken)
	user, err := middleware.AuthenticateAndFetchUser(r, models.UserTypeAppUser)
	if err != nil {
		networkutil.WSError(
			http.StatusUnauthorized,
			err.Error(),
		).Send(conn)
		return
	}

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

	roomManager.AddClient(&client, parsedRoomID)
	defer roomManager.RemoveClient(&client, parsedRoomID)

	for {
		_, rawMsg, err := conn.ReadMessage()
		if err != nil {
			logger.Log.Errorf("read error, %s", err)
			break
		}
		logger.Log.Infof("received in room %s: %s", parsedRoomID, string(rawMsg))
		roomManager.Broadcast(parsedRoomID, user.AppUser.UserID, string(rawMsg))
	}
}

func (m *RoomManager) Broadcast(roomID uuid.UUID, userID uuid.UUID, msg string) {
	fmt.Println("incoming: ", msg)
	m.mu.Lock()
	defer m.mu.Unlock()

	roomClients, exists := m.rooms[roomID]
	if !exists {
		return
	}

	// Persist only the content of the message
	savedMsg, err := communityservice.SaveCommunityMessage(
		context.Background(),
		roomID,
		userID,
		msg,
	)

	if err != nil {
		logger.Log.Errorf("failed to save message: %s", err)
		return
	}

	jsonMsg, err := json.Marshal(dto.SocketMessageDTO{
		MessageId:      savedMsg.ID,
		MessageGroupID: uuid.New(),
		Name:           savedMsg.Name.String,
		Username:       savedMsg.Username.String,
		UserProfileUrl: savedMsg.ProfilePictureUrl.String,
		UserColor:      savedMsg.Color,
		UserID:         savedMsg.UserID,
		CommunityID:    savedMsg.CommunityID,
		Content:        savedMsg.Content.String,
		Timestamp:      savedMsg.CreatedAt.Time,
	})

	if err != nil {
		logger.Log.Errorf("failed to marshal message: %s", err)
		return
	}

	for client := range roomClients {
		err := client.Conn.WriteMessage(websocket.TextMessage, jsonMsg)
		if err != nil {
			logger.Log.Errorf("write error: %s", err)
			client.Conn.Close()
			delete(roomClients, client)
		}
	}
}
