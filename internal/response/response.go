package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"strconv"
)

type StatusCode int

const (
	ResponseOK                  StatusCode = 200
	ResponseBadRequest          StatusCode = 400
	ResponseInternalServerError StatusCode = 500
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := headers.Headers{}
	header.Set("Content-Length", strconv.Itoa(contentLen))
	header.Set("Connection", "close")
	header.Set("Content-Type", "text/plain")
	return header
}

func getStatusLine(statusCode StatusCode) []byte {
	reasonPhrase := ""
	switch statusCode {
	case ResponseOK:
		reasonPhrase = "OK"
	case ResponseBadRequest:
		reasonPhrase = "Bad Request"
	case ResponseInternalServerError:
		reasonPhrase = "Internal Server Error"
	}
	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase))
}
