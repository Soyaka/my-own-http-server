package main

import (
	"fmt"
	"os"
	"github.com/codecrafters-io/http-server-starter-go/app/server"
)

const (
	HOST = "0.0.0.0"
	PORT = "4221"
)

func main() {
	fmt.Println("Logs from your program will appear here!")
	server := server.NewServer(HOST, PORT)
	err := server.Run()
	if err != nil {
		os.Exit(1)
	}
}
