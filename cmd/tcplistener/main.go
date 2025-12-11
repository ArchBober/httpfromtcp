package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
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

			request, err := request.RequestFromReader(c)
			if err != nil {
				log.Fatalf("could not resolve request: %v", err)
			}
			fmt.Printf(
				"Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n",
				request.RequestLine.Method,
				request.RequestLine.RequestTarget,
				request.RequestLine.HttpVersion,
			)
		}(conn)
	}
}
