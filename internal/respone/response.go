package respone

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	ResponseOK                  StatusCode = 200
	ResponseBadRequest          StatusCode = 400
	ResponseInternalServerError StatusCode = 500
)

const crlf = "\r\n"

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	switch statusCode {
	case ResponseOK:
		_, err := w.Write([]byte(fmt.Sprintf("HTTP/1.1 %v OK%s", statusCode, crlf)))
		if err != nil {
			return err
		}
	case ResponseBadRequest:
		_, err := w.Write([]byte(fmt.Sprintf("HTTP/1.1 %v Bad Request%s", statusCode, crlf)))
		if err != nil {
			return err
		}
	case ResponseInternalServerError:
		_, err := w.Write([]byte(fmt.Sprintf("HTTP/1.1 %v Internal Server Error%s", statusCode, crlf)))
		if err != nil {
			return err
		}
	default:
		_, err := w.Write([]byte(fmt.Sprintf("HTTP/1.1 %v %s", statusCode, crlf)))
		if err != nil {
			return err
		}
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := headers.Headers{}
	header.Set("Content-Length", strconv.Itoa(contentLen))
	header.Set("Connection", "close")
	header.Set("Content-Type", "text/plain")
	return header
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	contLen, ok := headers.Get("Content-Length")
	if !ok {
		return fmt.Errorf("Conent-Length header not declared")
	}

	conn, ok := headers.Get("Connection")
	if !ok {
		return fmt.Errorf("connection header not declared")
	}

	contType, ok := headers.Get("Content-Type")
	if !ok {
		return fmt.Errorf("Conent-Type header not declared")
	}

	_, err := w.Write([]byte(fmt.Sprintf("Content-Length: %s%s", contLen, crlf)))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(fmt.Sprintf("Connection: %s%s", conn, crlf)))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(fmt.Sprintf("Content-Type: %s%s", contType, crlf)))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(fmt.Sprint(crlf)))
	if err != nil {
		return err
	}

	return nil
}
