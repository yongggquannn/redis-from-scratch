package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Listening on port :6379")

	// Initialize a TCP listener on port 6379
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Listen for incoming connections
	connection, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()

	for {
		// Initialize a new RESP parser and read RESP value from connection
		resp := NewResp(connection)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(value)

		// ignore request and send back a PONG
		connection.Write([]byte("+OK\r\n"))
	}
	
}