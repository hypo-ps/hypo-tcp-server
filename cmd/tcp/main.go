package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"tcp/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("PORT???")
		os.Exit(1)
	}
	port, err := strconv.ParseUint(os.Args[1], 10, 16)
	if err != nil {
		fmt.Printf("Invalid Port %s!!\n", os.Args[1])
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	s := server.NewTCPServer(uint16(port))
	go func() {
		if err := s.Start(); err != nil {
			slog.Error("erro while starting the server", "err", err)
		}
		close(done)
	}()

	<-done
	s.ShutDown()
	slog.Info("server closed!")
}
