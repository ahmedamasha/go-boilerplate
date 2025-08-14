package services

import (
	"cusror_ai/internal/models"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WebSocketConnection represents a WebSocket connection
type WebSocketConnection struct {
	conn   *websocket.Conn
	send   chan []byte
	userID *string // Optional user ID for targeted messages
}

// WebSocketService manages WebSocket connections and broadcasts
type WebSocketService struct {
	// Hub for managing connections
	connections map[*WebSocketConnection]bool
	broadcast   chan []byte
	register    chan *WebSocketConnection
	unregister  chan *WebSocketConnection
	mutex       sync.RWMutex
	upgrader    websocket.Upgrader
}

// NewWebSocketService creates a new WebSocket service
func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		connections: make(map[*WebSocketConnection]bool),
		broadcast:   make(chan []byte),
		register:    make(chan *WebSocketConnection),
		unregister:  make(chan *WebSocketConnection),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for development
				// In production, implement proper origin checking
				return true
			},
		},
	}
}

// Run starts the WebSocket hub
func (ws *WebSocketService) Run() {
	for {
		select {
		case conn := <-ws.register:
			ws.mutex.Lock()
			ws.connections[conn] = true
			ws.mutex.Unlock()
			log.Printf("WebSocket client connected. Total connections: %d", len(ws.connections))

		case conn := <-ws.unregister:
			ws.mutex.Lock()
			if _, ok := ws.connections[conn]; ok {
				delete(ws.connections, conn)
				close(conn.send)
			}
			ws.mutex.Unlock()
			log.Printf("WebSocket client disconnected. Total connections: %d", len(ws.connections))

		case message := <-ws.broadcast:
			ws.mutex.RLock()
			for conn := range ws.connections {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(ws.connections, conn)
				}
			}
			ws.mutex.RUnlock()
		}
	}
}

// HandleWebSocket handles WebSocket connection requests
func (ws *WebSocketService) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Get user ID from query parameters or headers
	userID := r.URL.Query().Get("user_id")
	var userIDPtr *string
	if userID != "" {
		userIDPtr = &userID
	}

	wsConn := &WebSocketConnection{
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userIDPtr,
	}

	// Register the connection
	ws.register <- wsConn

	// Start goroutines for reading and writing
	go ws.writePump(wsConn)
	go ws.readPump(wsConn)
}

// writePump pumps messages from the hub to the WebSocket connection
func (ws *WebSocketService) writePump(wsConn *WebSocketConnection) {
	defer func() {
		wsConn.conn.Close()
	}()

	for {
		select {
		case message, ok := <-wsConn.send:
			if !ok {
				wsConn.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := wsConn.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		}
	}
}

// readPump pumps messages from the WebSocket connection to the hub
func (ws *WebSocketService) readPump(wsConn *WebSocketConnection) {
	defer func() {
		ws.unregister <- wsConn
		wsConn.conn.Close()
	}()

	// Set read deadline and pong handler for keepalive
	wsConn.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	wsConn.conn.SetPongHandler(func(string) error {
		wsConn.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := wsConn.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket unexpected close error: %v", err)
			}
			break
		}
	}
}

// BroadcastMessage broadcasts a message to all connected clients
func (ws *WebSocketService) BroadcastMessage(message *models.WebSocketMessage) error {
	data, err := message.ToJSON()
	if err != nil {
		return err
	}

	select {
	case ws.broadcast <- data:
	default:
		log.Printf("WebSocket broadcast channel full, dropping message")
	}

	return nil
}

// SendToUser sends a message to a specific user
func (ws *WebSocketService) SendToUser(userID string, message *models.WebSocketMessage) error {
	data, err := message.ToJSON()
	if err != nil {
		return err
	}

	ws.mutex.RLock()
	defer ws.mutex.RUnlock()

	sent := false
	for conn := range ws.connections {
		if conn.userID != nil && *conn.userID == userID {
			select {
			case conn.send <- data:
				sent = true
			default:
				// Channel full, close connection
				close(conn.send)
				delete(ws.connections, conn)
			}
		}
	}

	if !sent {
		log.Printf("No active WebSocket connection found for user %s", userID)
	}

	return nil
}

// Notification methods

// NotifyEventProcessed notifies about event processing
func (ws *WebSocketService) NotifyEventProcessed(event *models.Event, userID string) {
	message := models.NewEventProcessedMessage(event, userID)

	// Send to specific user if specified
	if userID != "" {
		ws.SendToUser(userID, message)
	} else {
		ws.BroadcastMessage(message)
	}
}

// NotifySegmentChanged notifies about segment changes
func (ws *WebSocketService) NotifySegmentChanged(userID string, oldSegments, newSegments []models.Segment) {
	message := models.NewSegmentChangedMessage(userID, oldSegments, newSegments)

	// Send to specific user
	ws.SendToUser(userID, message)
}

// NotifyOfferGenerated notifies about offer generation
func (ws *WebSocketService) NotifyOfferGenerated(offer *models.Offer, userID *string, segmentID uuid.UUID) {
	message := models.NewOfferGeneratedMessage(offer, userID, segmentID)

	// Send to specific user if specified, otherwise broadcast
	if userID != nil {
		ws.SendToUser(*userID, message)
	} else {
		ws.BroadcastMessage(message)
	}
}

// NotifyOfferUsed notifies about offer usage
func (ws *WebSocketService) NotifyOfferUsed(offer *models.Offer, userID string, orderID *string, amount float64) {
	message := models.NewOfferUsedMessage(offer, userID, orderID, amount)

	// Send to specific user
	ws.SendToUser(userID, message)
}

// NotifyError sends error notification
func (ws *WebSocketService) NotifyError(code, errorMessage, details string, userID *string) {
	message := models.NewErrorMessage(code, errorMessage, details, userID)

	// Send to specific user if specified, otherwise broadcast
	if userID != nil {
		ws.SendToUser(*userID, message)
	} else {
		ws.BroadcastMessage(message)
	}
}

// GetConnectionCount returns the number of active connections
func (ws *WebSocketService) GetConnectionCount() int {
	ws.mutex.RLock()
	defer ws.mutex.RUnlock()
	return len(ws.connections)
}

// GetConnectedUsers returns a list of connected user IDs
func (ws *WebSocketService) GetConnectedUsers() []string {
	ws.mutex.RLock()
	defer ws.mutex.RUnlock()

	userSet := make(map[string]bool)
	for conn := range ws.connections {
		if conn.userID != nil {
			userSet[*conn.userID] = true
		}
	}

	users := make([]string, 0, len(userSet))
	for userID := range userSet {
		users = append(users, userID)
	}

	return users
}

// SendHeartbeat sends heartbeat to maintain connections
func (ws *WebSocketService) SendHeartbeat() {
	heartbeatMessage := &models.WebSocketMessage{
		ID:        uuid.New(),
		Type:      "heartbeat",
		Data:      map[string]interface{}{"timestamp": time.Now()},
		Timestamp: time.Now(),
	}

	ws.BroadcastMessage(heartbeatMessage)
}

// StartHeartbeat starts periodic heartbeat
func (ws *WebSocketService) StartHeartbeat() {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			ws.SendHeartbeat()
		}
	}()
}
