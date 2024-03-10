package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	HOST          string = "0.0.0.0"
	PORT          string = "4221"
	ContentType   string = "Content-Type"
	ContentLength string = "Content-Length"
	Ok            string = "HTTP/1.1 200 OK"
	KO            string = "HTTP/1.1 404 Not Found"
	SEPARATOR     string = "\r\n"
	OCT           string = "application/octet-stream"
	TXT           string = "text/plain"
	SLASH         string = "/"
)

type server struct {
	Host string
	Port string
}

func NewServer(host, port string) *server {
	return &server{
		Host: host,
		Port: port,
	}
}

type file struct {
	Dir     string
	Name    string
	Content string
}

type request struct {
	Method  string
	Path    string
	Version string
	Body    string
	Conn    net.Conn
	Host    string
}

func main() {

	fmt.Println("Logs from your program will appear here!")
	server := NewServer(HOST, PORT)
	listner, err := net.Listen("tcp", server.Host+":"+server.Port)
	if err != nil {
		os.Exit(1)
	}
	for {
		conn, err := listner.Accept()
		if err != nil {
			fmt.Println("can,t accept incomming connections", err)
			continue
		}
		go server.Run(conn)
	}

}

func NewRequest(conn *net.Conn) *request {
	return &request{
		Conn: *conn,
	}
}

func (s *server) Run(conn net.Conn) error {
	defer conn.Close()
	request := NewRequest(&conn)
	err := request.requestParser()
	if err != nil {
		fmt.Println("cant initialize new request or send response", err)

		return err
	}
	if request.Body != "" {
		err = request.SendResponse()
		if err != nil {
			fmt.Println("cant send response")
			return err
		}
	}
	return nil
}

//file section

func NewFile() *file {
	return &file{}
}

func flagCheck() bool {
	if len(os.Args) > 2 {
		switch strings.ToLower(os.Args[1]) {
		case "--directory":
			return true
		}
	}
	return false
}

func SetFileDir() string {
	return os.Args[2]
}

func CheckFileExist(file *file) bool {
	_, err := os.Stat(file.Dir + string(os.PathSeparator) + file.Name)
	return err == nil
}

//end file section

func (r *request) requestParser() error {
	buff := make([]byte, 1024)
	_, err := r.Conn.Read(buff)
	if err != nil {
		return err
	}
	verbs := strings.Split(string(buff), "\r\n")

	if strings.Contains(verbs[0], "GET / HTTP/1.1") {
		r.Method = "GET"
		r.Path = "/"
		r.Version = "HTTP/1.1"
		r.Body = Ok + SEPARATOR + SEPARATOR
		return nil
	} else if strings.Contains(verbs[0], "GET /index.html HTTP/1.1") {
		r.Method = "GET"
		r.Path = "/index.html"
		r.Version = "HTTP/1.1"
		r.Body = KO + SEPARATOR + SEPARATOR
		return nil
	} else if strings.Contains(verbs[0], "GET /user-agent") {
		head := strings.Split(verbs[2], " ")
		if strings.Contains(head[0], "User-Agent:") {
			body := strings.Join(head[1:], "")
			if body != "" {
				r.Body = ResponseMaker(body, TXT)

			} else {
				r.Body = KO + SEPARATOR + SEPARATOR
			}
		} else {
			r.Body = KO + SEPARATOR + SEPARATOR
		}
	} else if strings.Contains(verbs[0], "GET /echo") {
		head := strings.Split(verbs[0], " ")
		r.Method = head[0]
		r.Path = head[1]
		r.Version = head[2]
		pathBody := strings.Split(head[1], "/")
		body := strings.Join(pathBody[2:], "/")
		r.Body = ResponseMaker(body, TXT)
	} else if strings.Contains(verbs[0], "GET /files/") {
		head := strings.Split(verbs[0], " ")
		r.Method = head[0]
		r.Path = head[1]
		r.Version = head[2]
		pathBody := strings.Split(head[1], "/")
		filePath := strings.Join(pathBody[2:], "/")
		if flagCheck() {
			// The if Triangle hell stars Here
			file := NewFile()
			file.Dir = SetFileDir()
			file.Name = filePath
			if CheckFileExist(file) {
				data, err := os.ReadFile(file.Dir + string(os.PathSeparator) + file.Name)
				if err != nil {
					r.Body = KO + SEPARATOR + SEPARATOR
				}
				r.Body = ResponseMaker(string(data), OCT)
				fmt.Println(r.Body)

			} else {
				r.Body = KO + SEPARATOR + SEPARATOR
			}
		} else {
			r.Body = KO + SEPARATOR + SEPARATOR
		}

	} else {
		r.Body = KO + SEPARATOR + SEPARATOR
	}
	return nil
}

func (r *request) SendResponse() error {
	_, err := r.Conn.Write([]byte(r.Body))
	if err != nil {
		return err
	}
	return nil
}

func ResponseMaker(respBody string, format string) string {
	return Ok + SEPARATOR + ContentType + ":" + format + SEPARATOR + ContentLength + ":" + fmt.Sprintf("%d", len(respBody)) + SEPARATOR + SEPARATOR + respBody
}
