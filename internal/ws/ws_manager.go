package ws

import (
	"encoding/json"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/websocket"
)

type Message struct {
	Type      string      `json:"type"`
	SessionID string      `json:"session_id,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type Session struct {
	Conn    *websocket.Conn
	mu      sync.Mutex
	LastMsg time.Time
}

type Manager struct {
	sessions sync.Map // sessionID -> *Session
	connCount atomic.Int64
}

var ManagerInstance = &Manager{}

func (m *Manager) RegisterSession(sessionID string, conn *websocket.Conn) {
	m.sessions.Store(sessionID, &Session{Conn: conn, LastMsg: time.Now()})
	m.connCount.Add(1)
	log.Printf("WebSocket session registered: %s (total: %d)", sessionID, m.connCount.Load())
	
	// Send welcome message
	m.sendToSession(sessionID, Message{
		Type:      "connected",
		SessionID: sessionID,
		Data:      map[string]string{"message": "Connected to progress stream"},
	})
}

func (m *Manager) UnregisterSession(sessionID string) {
	val, ok := m.sessions.LoadAndDelete(sessionID)
	if !ok {
		return
	}
	if s, ok := val.(*Session); ok {
		s.Conn.Close()
	}
	m.connCount.Add(-1)
	log.Printf("WebSocket session removed: %s (total: %d)", sessionID, m.connCount.Load())
}

func (m *Manager) SendProgress(sessionID string, progress float64, status string, message string) {
	m.sendToSession(sessionID, Message{
		Type:      "progress",
		SessionID: sessionID,
		Data: map[string]interface{}{
			"progress": progress,
			"status":   status,
			"message":  message,
		},
	})
}

func (m *Manager) SendComplete(sessionID string, data interface{}) {
	m.sendToSession(sessionID, Message{
		Type:      "complete",
		SessionID: sessionID,
		Data:      data,
	})
}

func (m *Manager) SendError(sessionID string, err string) {
	m.sendToSession(sessionID, Message{
		Type:      "error",
		SessionID: sessionID,
		Data:      map[string]string{"error": err},
	})
}

func (m *Manager) sendToSession(sessionID string, msg Message) {
	val, ok := m.sessions.Load(sessionID)
	if !ok {
		return
	}
	s := val.(*Session)
	s.mu.Lock()
	defer s.mu.Unlock()
	
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal WS message: %v", err)
		return
	}
	
	if err := websocket.Message.Send(s.Conn, string(data)); err != nil {
		log.Printf("Failed to send WS message: %v", err)
		m.sessions.Delete(sessionID)
		m.connCount.Add(-1)
	} else {
		s.LastMsg = time.Now()
	}
}

func (m *Manager) Handler() websocket.Handler {
	return websocket.Handler(func(conn *websocket.Conn) {
		sessionID := conn.Request().URL.Query().Get("session")
		if sessionID == "" {
			sessionID = conn.RemoteAddr().String()
		}
		
		m.RegisterSession(sessionID, conn)
		defer m.UnregisterSession(sessionID)
		
		// Keep connection alive and handle messages
		for {
			var msg string
			if err := websocket.Message.Receive(conn, &msg); err != nil {
				log.Printf("WS read error: %v", err)
				break
			}
			
			// Parse ping messages
			var m Message
			if json.Unmarshal([]byte(msg), &m) == nil && m.Type == "ping" {
				m.Type = "pong"
				data, _ := json.Marshal(m)
				websocket.Message.Send(conn, string(data))
			}
		}
	})
}

func (m *Manager) Stats() map[string]interface{} {
	count := m.connCount.Load()
	return map[string]interface{}{
		"active_sessions": count,
		"total_connections": count,
	}
}
