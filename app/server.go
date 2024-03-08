package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	RESPONSEOK = "HTTP/1.1 200 OK\r\n\r\n"
	RESPONSEKO = "HTTP/1.1 404 Not Found\r\n\r\n"

)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	Conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	buff := make([]byte, 1024)
	_, err = Conn.Read(buff)
	if err!= nil{
		os.Exit(1)
	}
	if strings.Contains(string(buff), "GET / HTTP/1.1"){
		Conn.Write([]byte(RESPONSEOK))
	}else{
		Conn.Write([]byte(RESPONSEKO))
	}

}
