package manager

import (
	"fmt"
	"gobloks/internal/types"
	"gobloks/internal/utilities"
	"sync"

	"github.com/gorilla/websocket"
)

type SocketConnection struct {
	socket *websocket.Conn
	mu     *sync.Mutex
}

func InitSocketConnection(socket *websocket.Conn) *SocketConnection {
	return &SocketConnection{socket, &sync.Mutex{}}
}

type SocketManager struct {
	activeConnections utilities.Set[*SocketConnection]
	mu                *sync.Mutex
}

func InitSocketManager() *SocketManager {
	return &SocketManager{utilities.NewSet([]*SocketConnection{}), &sync.Mutex{}}
}

func (s *SocketManager) Connect(socket *websocket.Conn) *SocketConnection {
	s.mu.Lock()
	defer s.mu.Unlock()
	conn := InitSocketConnection(socket)
	s.activeConnections.Add(conn)
	return conn
}

func (s *SocketManager) Disconnect(conn *SocketConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.activeConnections.Remove(conn)
}

func (s *SocketManager) Send(conn *SocketConnection, out *types.SocketData) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	err := conn.socket.WriteJSON(out)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *SocketManager) Recv(conn *SocketConnection, in *types.SocketData) error {
	return conn.socket.ReadJSON(in)
}

// func (s *SocketManager) Broadcast(message *types.ChatMessage) {
func (s *SocketManager) Broadcast(out *types.SocketData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for socket := range s.activeConnections {
		go s.Send(socket, out)
	}
}
