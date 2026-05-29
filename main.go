package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer f.Close()
		defer close(lines)

		buffer := make([]byte, 8)
		currentLine := ""

		for {
			n, err := f.Read(buffer)

			if err == io.EOF {
				break
			}

			if err != nil {
				return
			}

			parts := strings.Split(string(buffer[:n]), "\n")

			for i := 0; i < len(parts)-1; i++ {
				lines <- currentLine + parts[i]
				currentLine = ""
			}

			currentLine += parts[len(parts)-1]
		}

		if currentLine != "" {
			lines <- currentLine
		}
	}()

	return lines
}

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	for line := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", line)
	}
}