package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, err
	}
	headerData := string(data)
	requestLine, err := parseRequestLine(headerData)
	if err != nil {
		fmt.Printf("Error parsing header data: %v", err)
		return &Request{}, err
	}

	request := Request{RequestLine: requestLine}
	return &request, nil
}

func parseRequestLine(line string) (RequestLine, error) {
	requestLine := strings.Split(line, "\r\n")[0]
	if len(requestLine) < 1 {
		return RequestLine{}, fmt.Errorf("got empty header")
	}

	splitRequestLine := strings.Split(requestLine, " ")
	if len(splitRequestLine) != 3 {
		return RequestLine{}, fmt.Errorf("got invalid request line")
	}

	requestTarget := splitRequestLine[1]
	if strings.Split(requestTarget, "")[0] != "/" {
		return RequestLine{}, fmt.Errorf("got wrong path not starting with '/'")
	}

	method := splitRequestLine[0]
	allowed := []string{"GET", "POST"}
	isAllowed := false
	for _, m := range allowed {
		if method == m {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return RequestLine{}, fmt.Errorf("not allowed method: %s", method)
	}

	splitHttp := strings.Split(splitRequestLine[2], "/")
	if splitHttp[0] != "HTTP" || splitHttp[1] != "1.1" {
		return RequestLine{}, fmt.Errorf("not supported http version: %s", splitRequestLine[2])
	}

	return RequestLine{
		HttpVersion:   splitHttp[1],
		RequestTarget: requestTarget,
		Method:        method,
	}, nil
}
