package main

import (
	"fmt"
	"io"
	"net"
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
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("connection accepted")

		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}

		fmt.Println("connection closed")
	}
}