package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("Could not open file: %v", err)
		os.Exit(1)
	}
	defer f.Close()

	buff := make([]byte, 8)
	var current_line string
	for {
		n, err := f.Read(buff)
		splitted_string := strings.Split(string(buff[:n]), "\n")
		current_line = current_line + splitted_string[0]
		if len(splitted_string) != 1 {
			fmt.Printf("read: %s\n", current_line)
			current_line = splitted_string[1]
		}
		if err != nil {
			if err == io.EOF {
				os.Exit(0)
			}
			log.Fatal("Could not read file: %v", err)
			os.Exit(1)
		}
	}

}
