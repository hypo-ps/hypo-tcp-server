package server

import (
	"io"
	"log/slog"
	"net"
)

type Connection struct {
	conn          net.Conn
	clientDetails net.Addr
}

func NewConnection(listner net.Listener) (*Connection, error) {
	conn, err := listner.Accept()
	if err != nil {
		slog.Error("error while accepting connection", "err", err)
		return nil, err
	}

	return &Connection{
		conn:          conn,
		clientDetails: conn.RemoteAddr(),
	}, nil
}

func (c *Connection) HandleConnection() {
	packet := make([]byte, 1024)
	defer c.Close()
	for {
		n, err := c.conn.Read(packet)
		if err != nil {
			if err == io.EOF {
				return
			}
			slog.Error("error while reading the packet", "err", err)
		}

		slog.Info("message received: ", "packet", string(packet[:n]))
		c.conn.Write([]byte("Message received! Thank you!"))
	}
}

func (c *Connection) Close() {
	slog.Info("Clossing conenction with ", "client", c.clientDetails.String())
	c.conn.Close()
}
