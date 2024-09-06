package sockets

import (
	"fmt"
	"gobloks/internal/types"
	"gobloks/internal/utilities"
	"sync"

	"github.com/gorilla/websocket"
)

// Message types for socket communication
const (
	PLAYER_UPDATE types.SocketDataType = iota
	BOARD_STATE
	PRIVATE_GAME_STATE
	CHAT_MESSAGE
	GAME_STATUS
)

type Connection struct {
	socket *websocket.Conn
	mu     *sync.Mutex
}

func initConnection(socket *websocket.Conn) *Connection {
	return &Connection{socket, &sync.Mutex{}}
}

func (s *Connection) send(out *types.SocketData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.socket.WriteJSON(out)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *Connection) recv(in *types.SocketData) error {
	return s.socket.ReadJSON(in)
}

type SocketManager struct {
	activeConnections utilities.Set[*Connection]
	mu                *sync.Mutex
}

func InitSocketManager(size int) *SocketManager {
	return &SocketManager{utilities.NewSet([]*Connection{}, size), &sync.Mutex{}}
}

func (s *SocketManager) Connect(ws *websocket.Conn) *Connection {
	s.mu.Lock()
	defer s.mu.Unlock()
	conn := initConnection(ws)
	s.activeConnections.Add(conn)
	return conn
}

func (s *SocketManager) Disconnect(conn *Connection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.activeConnections.Remove(conn)
}

func (s *SocketManager) Send(conn *Connection, out *types.SocketData) {
	go conn.send(out)
}

func (s *SocketManager) Recv(conn *Connection, in *types.SocketData) error {
	return conn.recv(in)
}

// func (s *SocketManager) Broadcast(message *types.ChatMessage) {
func (s *SocketManager) Broadcast(out *types.SocketData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for socket := range s.activeConnections {
		go socket.send(out)
	}
}
