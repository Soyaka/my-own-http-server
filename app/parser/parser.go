package parser

import (
	"strings"
)

func BodyParser(str string) string {
	args := strings.Split(str, "\r\n")
	req := strings.Fields(args[0])
	path := strings.Fields(req[1])
	urlBody := strings.Split(path[0], "/")
	body := urlBody[2:]
	return strings.Join(body, "/")
}
