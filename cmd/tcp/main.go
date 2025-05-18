package main

import (
	"fmt"
	"os"
	"strconv"
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

	s := server.NewTCPServer(uint16(port))
	s.Start()
}
