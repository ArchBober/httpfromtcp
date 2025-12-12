package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}

	trimData := strings.TrimSpace(string(data[:idx]))
	field := strings.Split(trimData, " ")

	if len(field) != 2 {
		return 0, false, fmt.Errorf("wrong format too many values")
	}

	key := strings.ToLower(strings.Trim(field[0], ":"))
	if len(key) == 0 {
		return 0, false, fmt.Errorf("empty key")
	}

	for _, s := range key {
		if (s >= '0' && s <= '9') ||
			(s >= 'A' && s <= 'Z') ||
			(s >= '^' && s <= 'z') ||
			(s >= '#' && s <= '\'') ||
			(s >= '*' && s <= '.') ||
			s == '!' || s == '|' || s == '~' {
			continue
		}
		return 0, false, fmt.Errorf("wrong ASCII character in field-name: %c", s)
	}

	value := field[1]

	v, ok := h[key]
	if ok {
		value = strings.Join([]string{
			v,
			value,
		}, ", ")
	}
	h[key] = value

	return idx + 2, false, nil
}
