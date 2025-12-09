package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const adresstcp = ":42069"

func main() {
	listener, err := net.Listen("tcp", adresstcp)
	if err != nil {
		log.Fatalf("could not create listener: %s\n", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Could not accept connection: %s", err)
		}
		fmt.Println("Connection Accepted")

		go func(c net.Conn) {
			defer c.Close()
			defer fmt.Println("Connection Closed")

			strCh := getLinesChannel(c)
			for line := range strCh {
				fmt.Println(line)
			}
		}(conn)
	}
}

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
