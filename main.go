package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	fmt.Println("Listening on port :6379")

	// Initialize a TCP listener on port 6379
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Listen for incoming connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	// Create infinite loop to read from the connection
	for {
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client")
				break
			}
			fmt.Println("Error reading from client:", err.Error())
			os.Exit(1)

		}
		// Ignore request and send back a PONG
		conn.Write([]byte("+OK\r\n"))
	}
	
}