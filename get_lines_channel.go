package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	stringCh := make(chan string)
	currentLineContents := ""

	go func() {
		defer f.Close()
		defer close(stringCh)
		for {
			buff := make([]byte, 8, 8)
			n, err := f.Read(buff)
			if err != nil {
				if currentLineContents != "" {
					stringCh <- currentLineContents
					currentLineContents = ""
				}
				if errors.Is(err, io.EOF) {
					return
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			str := string(buff[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				currentLineContents = currentLineContents + parts[i]
				stringCh <- currentLineContents
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return stringCh
}
