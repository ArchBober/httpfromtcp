package main

import (
	"fmt"
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

			strCh := getLinesChannel(c)
			for line := range strCh {
				fmt.Println(line)
			}
		}(conn)
	}
}
