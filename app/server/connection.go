package server

import (
	"fmt"
	"net"
	"strconv"

	"github.com/codecrafters-io/http-server-starter-go/app/parser"
)

const (
	RESPONSEOK    = "HTTP/1.1 200 OK"
	RESPONSEKO    = "HTTP/1.1 404 Not Found"
	SEPARATOR     = "\r\n"
	ContentType   = "Content-Type"
	ContentLength = "Content-Length"
	TXT           = "text/plain"
	SEP           = ":"
)

type Server struct {
	Host     string
	Port     string
	Protocol string
	Conn     net.Conn
}

func NewServer(host string, port string) *Server {
	return &Server{
		Host:     host,
		Port:     port,
		Protocol: "tcp",
	}
}

func (s *Server) Run() error {
	Listener, err := net.Listen(s.Protocol, s.Host+":"+s.Port)
	if err != nil {
		fmt.Println("can't start server", err)
		return err

	}
	conn, err := Listener.Accept()
	if err != nil {
		fmt.Println("cant't accept connection", err)
		return err
	}
	s.Conn = conn
	buff := make([]byte, 1024)
	if _, err := s.Conn.Read(buff); err != nil {
		fmt.Print("can't read from buffer", err)
		return err
	}

	err = s.ResponseWriter(buff)
	if err != nil {
		fmt.Print("can't Write Response")
		return err

	}
	return nil

}

func (s *Server) ResponseWriter(buff []byte) error {
	var response string
	body := parser.BodyParser(string(buff))
	if body != "" {
		response = RESPONSEOK + SEPARATOR + ContentType + SEP + TXT + SEPARATOR + ContentLength + ":" + strconv.Itoa(len(body)) + SEPARATOR + SEPARATOR + body
		_, err := s.Conn.Write([]byte(response))

		if err != nil {
			return err
		}
	} else {
		response = RESPONSEKO + SEPARATOR + SEPARATOR
		_, err := s.Conn.Write([]byte(response))
		if err != nil {
			return err
		}
	}

	return nil
}
