package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state       int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buff := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	request := Request{
		RequestLine: RequestLine{},
		state:       0,
	}

	for request.state != 1 {
		nBytes, err := reader.Read(buff[readToIndex:])
		if err != nil {
			if err == io.EOF {
				request.state = 1
				break
			}
			return nil, err
		}

		readToIndex += nBytes

		if readToIndex == len(buff) {
			newBuff := make([]byte, 2*len(buff))
			copy(newBuff, buff)
			buff = newBuff
		}

		parsednbytes, err := request.parse(buff[0:readToIndex])
		if err != nil {
			return &Request{}, err
		}

		copy(buff, buff[parsednbytes:readToIndex])
		readToIndex -= parsednbytes

	}

	return &request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return &RequestLine{}, 0, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.state == 0 {
		requestLine, nbytes, err := parseRequestLine(data)
		if err != nil {
			return nbytes, err
		}
		if nbytes == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = 1
		return nbytes, nil
	}
	if r.state == 1 {
		return 0, fmt.Errorf("error: trying to read data in a done state")
	}
	return 0, fmt.Errorf("error: unknown state")
}
