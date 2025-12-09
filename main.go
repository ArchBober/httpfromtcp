package main

import (
	"fmt"
	"log"
	"os"
)

const inputFilePath = "messages.txt"

func main() {
	f, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("could not open %s: %s\n", inputFilePath, err)
	}

	fmt.Printf("Reading data from %s\n", inputFilePath)
	fmt.Println("=====================================")

	strCh := getLinesChannel(f)
	for line := range strCh {
		fmt.Printf("read: %s\n", line)
	}
}
