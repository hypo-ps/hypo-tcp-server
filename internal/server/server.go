package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"
)

type TCPServer struct {
	listner        net.Listener
	connections    map[string]*Connection
	port           uint16
	lock           sync.Mutex
	connWg         *sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
	closedConnChan chan *Connection
}

func NewTCPServer(port uint16) *TCPServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &TCPServer{
		connections:    make(map[string]*Connection),
		port:           port,
		connWg:         &sync.WaitGroup{},
		ctx:            ctx,
		cancel:         cancel,
		closedConnChan: make(chan *Connection, 1),
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

	go s.HandleClosedConnection()
	go s.LogConnections()
	for {
		welcomeConnection(s)
	}
}

func welcomeConnection(s *TCPServer) error {
	c, err := NewConnection(&s.ctx, s.listner, s.closedConnChan)
	if err != nil {
		return err
	}
	slog.Info("new connection accepted from", "client", c.clientDetails)

	s.lock.Lock()
	s.connections[c.GetId()] = c
	s.lock.Unlock()

	s.connWg.Add(1)
	go c.Handle(s.connWg)
	return nil
}

func (s *TCPServer) ShutDown() {
	slog.Info("Shutting down...")
	s.cancel()
	s.connWg.Wait()
}

func (s *TCPServer) LogConnections() {
	for {
		time.Sleep(2 * time.Second)
		slog.Info("total active", "connections", len(s.connections))
	}
}

func (s *TCPServer) HandleClosedConnection() {
	for conn := range s.closedConnChan {
		slog.Info("handling closed connection")
		s.lock.Lock()
		delete(s.connections, conn.GetId())
		s.lock.Unlock()
	}
}
