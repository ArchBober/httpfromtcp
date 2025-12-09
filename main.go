package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("Could not open file: %v", err)
		os.Exit(1)
	}
	defer f.Close()

	buff := make([]byte, 8)

	for {
		n, err := f.Read(buff)
		if err != nil {
			if err == io.EOF {
				os.Exit(0)
			}
			log.Fatal("Could not read file: %v", err)
			os.Exit(1)
		}
		fmt.Printf("read: %s\n", string(buff[:n]))
	}

}
