package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/conneroisu/twerge/examples/realtime-collab/types"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin in development
	},
}

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	Broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Room management
	rooms map[string]*Room
	mu    sync.RWMutex

	// User presence tracking
	userPresence map[string]*UserPresence
	presenceMu   sync.RWMutex
}

// Room represents a collaboration room
type Room struct {
	ID      string
	Type    string // "project", "document", "general"
	Clients map[*Client]bool
	mu      sync.RWMutex
}

// UserPresence tracks real-time user information
type UserPresence struct {
	User      types.User
	LastSeen  time.Time
	IsOnline  bool
	CurrentRoom string
	Cursor    *types.Cursor
	Status    types.UserStatus
}

// Client represents a websocket connection
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	user     types.User
	rooms    map[string]bool
	lastPing time.Time
	mu       sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		Broadcast:    make(chan []byte),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		clients:      make(map[*Client]bool),
		rooms:        make(map[string]*Room),
		userPresence: make(map[string]*UserPresence),
	}
}

// Run starts the hub and handles client management
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.updateUserPresence(client.user.ID, client.user, true, "")
			
			// Send welcome message
			welcomeMsg := types.WebSocketMessage{
				Type:      types.MessageTypeUserJoined,
				UserID:    client.user.ID,
				Data:      client.user,
				Timestamp: time.Now(),
				MessageID: generateMessageID(),
			}
			
			if data, err := json.Marshal(welcomeMsg); err == nil {
				h.broadcastToAll(data)
			}

			log.Printf("Client registered: %s", client.user.Username)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.updateUserPresence(client.user.ID, client.user, false, "")
				h.leaveAllRooms(client)
				
				// Send user left message
				leftMsg := types.WebSocketMessage{
					Type:      types.MessageTypeUserLeft,
					UserID:    client.user.ID,
					Data:      client.user,
					Timestamp: time.Now(),
					MessageID: generateMessageID(),
				}
				
				if data, err := json.Marshal(leftMsg); err == nil {
					h.broadcastToAll(data)
				}

				log.Printf("Client unregistered: %s", client.user.Username)
			}

		case message := <-h.Broadcast:
			h.broadcastToAll(message)
		}
	}
}

// HandleWebSocket handles new websocket connections
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request, user types.User) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:      h,
		conn:     conn,
		send:     make(chan []byte, 256),
		user:     user,
		rooms:    make(map[string]bool),
		lastPing: time.Now(),
	}

	client.hub.register <- client

	// Start goroutines for this client
	go client.writePump()
	go client.readPump()
}

// Client read pump
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.lastPing = time.Now()
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var wsMsg types.WebSocketMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		c.handleMessage(wsMsg)
	}
}

// Client write pump
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Handle incoming messages from clients
func (c *Client) handleMessage(msg types.WebSocketMessage) {
	msg.UserID = c.user.ID
	msg.Timestamp = time.Now()
	msg.MessageID = generateMessageID()

	switch msg.Type {
	case types.MessageTypeJoinRoom:
		if roomID, ok := msg.Data.(string); ok {
			c.joinRoom(roomID, "general")
		}

	case types.MessageTypeLeaveRoom:
		if roomID, ok := msg.Data.(string); ok {
			c.leaveRoom(roomID)
		}

	case types.MessageTypeCursorMove:
		if cursorData, ok := msg.Data.(map[string]interface{}); ok {
			cursor := &types.Cursor{
				UserID:     c.user.ID,
				DocumentID: getString(cursorData, "documentId"),
				Line:       getInt(cursorData, "line"),
				Column:     getInt(cursorData, "column"),
				Color:      getString(cursorData, "color"),
				UpdatedAt:  time.Now(),
			}
			
			if selectionData, ok := cursorData["selection"].(map[string]interface{}); ok {
				cursor.Selection = &types.Selection{
					StartLine:   getInt(selectionData, "startLine"),
					StartColumn: getInt(selectionData, "startColumn"),
					EndLine:     getInt(selectionData, "endLine"),
					EndColumn:   getInt(selectionData, "endColumn"),
				}
			}
			
			c.hub.updateUserCursor(c.user.ID, cursor)
			c.hub.BroadcastToRoom(msg.RoomID, msg)
		}

	case types.MessageTypeUserTyping:
		c.hub.BroadcastToRoom(msg.RoomID, msg)

	case types.MessageTypeDocumentChange:
		c.hub.BroadcastToRoom(msg.RoomID, msg)

	case types.MessageTypeCommentAdded:
		c.hub.BroadcastToRoom(msg.RoomID, msg)

	case types.MessageTypePresenceUpdate:
		if statusData, ok := msg.Data.(map[string]interface{}); ok {
			if status, ok := statusData["status"].(string); ok {
				c.hub.updateUserStatus(c.user.ID, types.UserStatus(status))
				c.hub.broadcastToAll(marshalMessage(msg))
			}
		}

	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// Room management methods
func (c *Client) joinRoom(roomID, roomType string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.hub.mu.Lock()
	defer c.hub.mu.Unlock()

	// Create room if it doesn't exist
	if _, exists := c.hub.rooms[roomID]; !exists {
		c.hub.rooms[roomID] = &Room{
			ID:      roomID,
			Type:    roomType,
			Clients: make(map[*Client]bool),
		}
	}

	room := c.hub.rooms[roomID]
	room.mu.Lock()
	room.Clients[c] = true
	room.mu.Unlock()

	c.rooms[roomID] = true
	c.hub.updateUserPresence(c.user.ID, c.user, true, roomID)

	// Notify others in the room
	joinMsg := types.WebSocketMessage{
		Type:      types.MessageTypeUserJoined,
		RoomID:    roomID,
		UserID:    c.user.ID,
		Data:      c.user,
		Timestamp: time.Now(),
		MessageID: generateMessageID(),
	}

	c.hub.BroadcastToRoom(roomID, joinMsg)
}

func (c *Client) leaveRoom(roomID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.hub.mu.Lock()
	defer c.hub.mu.Unlock()

	if room, exists := c.hub.rooms[roomID]; exists {
		room.mu.Lock()
		delete(room.Clients, c)
		room.mu.Unlock()

		delete(c.rooms, roomID)

		// Clean up empty rooms
		if len(room.Clients) == 0 {
			delete(c.hub.rooms, roomID)
		}
	}

	// Notify others in the room
	leaveMsg := types.WebSocketMessage{
		Type:      types.MessageTypeUserLeft,
		RoomID:    roomID,
		UserID:    c.user.ID,
		Data:      c.user,
		Timestamp: time.Now(),
		MessageID: generateMessageID(),
	}

	c.hub.BroadcastToRoom(roomID, leaveMsg)
}

func (h *Hub) leaveAllRooms(client *Client) {
	client.mu.Lock()
	rooms := make([]string, 0, len(client.rooms))
	for roomID := range client.rooms {
		rooms = append(rooms, roomID)
	}
	client.mu.Unlock()

	for _, roomID := range rooms {
		client.leaveRoom(roomID)
	}
}

// Broadcasting methods
func (h *Hub) broadcastToAll(message []byte) {
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}

func (h *Hub) BroadcastToRoom(roomID string, msg types.WebSocketMessage) {
	h.mu.RLock()
	room, exists := h.rooms[roomID]
	h.mu.RUnlock()

	if !exists {
		return
	}

	data := marshalMessage(msg)
	
	room.mu.RLock()
	for client := range room.Clients {
		select {
		case client.send <- data:
		default:
			close(client.send)
			delete(room.Clients, client)
			delete(h.clients, client)
		}
	}
	room.mu.RUnlock()
}

// Presence management
func (h *Hub) updateUserPresence(userID string, user types.User, isOnline bool, currentRoom string) {
	h.presenceMu.Lock()
	defer h.presenceMu.Unlock()

	h.userPresence[userID] = &UserPresence{
		User:        user,
		LastSeen:    time.Now(),
		IsOnline:    isOnline,
		CurrentRoom: currentRoom,
		Status:      user.Status,
	}
}

func (h *Hub) updateUserCursor(userID string, cursor *types.Cursor) {
	h.presenceMu.Lock()
	defer h.presenceMu.Unlock()

	if presence, exists := h.userPresence[userID]; exists {
		presence.Cursor = cursor
		presence.LastSeen = time.Now()
	}
}

func (h *Hub) updateUserStatus(userID string, status types.UserStatus) {
	h.presenceMu.Lock()
	defer h.presenceMu.Unlock()

	if presence, exists := h.userPresence[userID]; exists {
		presence.Status = status
		presence.LastSeen = time.Now()
	}
}

// GetOnlineUsers returns all currently online users
func (h *Hub) GetOnlineUsers() []types.User {
	h.presenceMu.RLock()
	defer h.presenceMu.RUnlock()

	var users []types.User
	for _, presence := range h.userPresence {
		if presence.IsOnline {
			users = append(users, presence.User)
		}
	}
	return users
}

// GetRoomUsers returns all users in a specific room
func (h *Hub) GetRoomUsers(roomID string) []types.User {
	h.mu.RLock()
	room, exists := h.rooms[roomID]
	h.mu.RUnlock()

	if !exists {
		return []types.User{}
	}

	var users []types.User
	room.mu.RLock()
	for client := range room.Clients {
		users = append(users, client.user)
	}
	room.mu.RUnlock()

	return users
}

// GetUserCursors returns all active cursors in a document
func (h *Hub) GetUserCursors(documentID string) map[string]*types.Cursor {
	h.presenceMu.RLock()
	defer h.presenceMu.RUnlock()

	cursors := make(map[string]*types.Cursor)
	for userID, presence := range h.userPresence {
		if presence.IsOnline && presence.Cursor != nil && presence.Cursor.DocumentID == documentID {
			cursors[userID] = presence.Cursor
		}
	}
	return cursors
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID string, msg types.WebSocketMessage) {
	data := marshalMessage(msg)
	
	for client := range h.clients {
		if client.user.ID == userID {
			select {
			case client.send <- data:
			default:
				close(client.send)
				delete(h.clients, client)
			}
			break
		}
	}
}

// Utility functions
func marshalMessage(msg types.WebSocketMessage) []byte {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return []byte{}
	}
	return data
}

func generateMessageID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func getInt(data map[string]interface{}, key string) int {
	if val, ok := data[key].(float64); ok {
		return int(val)
	}
	return 0
}

// Health check for the hub
func (h *Hub) GetStats() map[string]interface{} {
	h.mu.RLock()
	h.presenceMu.RLock()
	defer h.mu.RUnlock()
	defer h.presenceMu.RUnlock()

	onlineCount := 0
	for _, presence := range h.userPresence {
		if presence.IsOnline {
			onlineCount++
		}
	}

	return map[string]interface{}{
		"total_clients":    len(h.clients),
		"total_rooms":      len(h.rooms),
		"online_users":     onlineCount,
		"total_users":      len(h.userPresence),
		"last_updated":     time.Now(),
	}
}