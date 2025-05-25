package server

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"
)

type Connection struct {
	id             string
	conn           net.Conn
	clientDetails  net.Addr
	ctx            context.Context
	cancel         context.CancelFunc
	closeOnce      sync.Once
	closedConnChan chan *Connection
}

func NewConnection(ctx *context.Context, listner net.Listener, closedConnChan chan *Connection) (*Connection, error) {
	conn, err := listner.Accept()
	if err != nil {
		slog.Error("error while accepting connection", "err", err)
		return nil, err
	}

	childCtx, cancel := context.WithCancel(*ctx)

	return &Connection{
		id:             generateID(),
		conn:           conn,
		clientDetails:  conn.RemoteAddr(),
		ctx:            childCtx,
		cancel:         cancel,
		closedConnChan: closedConnChan,
	}, nil
}

func (c *Connection) GetId() string {
	return c.id
}

func (c *Connection) Handle(wg *sync.WaitGroup) {
	packet := make([]byte, 1024)
	defer c.Close()
	defer wg.Done()
	go c.HandleCancel()
	for {
		select {
		case <-c.ctx.Done():
			slog.Info("context cancelled, closing connection with client", "client", c.clientDetails)
			return
		default:
			n, err := c.conn.Read(packet)
			if err != nil {
				if err == io.EOF {
					slog.Info("connection closed by", "client", c.clientDetails)
					return
				}

				if c.ctx.Err() != nil {
					slog.Info("context canceled, stoping handler", "client", c.clientDetails)
					return
				}

				slog.Error("error while reading the packet", "err", err)
			}

			slog.Info("message received: ", "packet", string(packet[:n]))
			c.conn.Write([]byte("Message received! Thank you!"))
		}
	}
}

func (c *Connection) Close() {
	c.closeOnce.Do(func() {
		slog.Info("Clossing conenction with", "client", c.clientDetails.String())
		err := c.conn.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("error while closing the connection with client %s", c.clientDetails), "err", err)
		}
		c.closedConnChan <- c
	})
}

func (c *Connection) HandleCancel() {
	<-c.ctx.Done()
	c.Close()
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
