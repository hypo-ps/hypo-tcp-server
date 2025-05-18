package server

import (
	"fmt"
	"log/slog"
	"net"
)

type TCPServer struct {
	listner     net.Listener
	connections []*Connection
	port        uint16
}

func NewTCPServer(port uint16) *TCPServer {
	return &TCPServer{
		connections: make([]*Connection, 0, 10),
		port:        port,
	}
}

func (s *TCPServer) Start() error {
	listner, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", s.port))
	if err != nil {
		slog.Error("error while starting to listen", "err", err)
		return fmt.Errorf("error while starting to listen: %w", err)
	}

	slog.Info("TCP TCPServer listning on port", "port", s.port)

	s.listner = listner

	for {
		s.WelcomeConnection()
	}
}

func (s *TCPServer) WelcomeConnection() error {
	c, err := NewConnection(s.listner)
	slog.Info("new connection accepted from: ", "client", c.clientDetails)
	if err != nil {
		return err
	}
	s.connections = append(s.connections, c)
	go c.HandleConnection()
	return nil
}
