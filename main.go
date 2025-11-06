package main

import (
	"fmt"
	"net"
	"strings"
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

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		// Extract command and perform validations to handle arguments
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(connection)

		handler, ok := Handlers[command]

		if !ok {
			fmt.Println("Invalid Command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		res := handler(args)
		writer.Write(res)
	}
	
}