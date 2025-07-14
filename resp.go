package main

import (
	"bufio"
	"io"
	"strconv"
	"fmt"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// Easier to parse & deserialize RESP values
type Value struct {
	typ   string // Determine data type of the value
	bulk  string // Bulk string value
	array []Value // Array of RESP values
}

type Resp struct {
	reader *bufio.Reader
}

// Pass buffer from the connection to the RESP parser
func NewResp(rd io.Reader) *Resp {
	return &Resp{
		reader: bufio.NewReader(rd),
	}
}

// Reads a line from the buffer
func (r *Resp) readLine() (line []byte, n int, err error) {
	line, err = r.reader.ReadBytes('\n')
	if err != nil {
		return nil, 0, err
	}
	n = len(line)
	// Remove CRLF (\r\n)
	if len(line) >= 2 && line[len(line)-2] == '\r' {
		return line[:len(line)-2], n, nil
	}
	return line, n, nil
}

// Reads an integer from the buffer
func (r *Resp) readInt() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	// Parse the integer value
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return int(i64), n, nil
}

// Generic function to read RESP values
func (r *Resp) Read() (Value, error) {
	// Read the first byte to determine the type of RESP value
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, fmt.Errorf("unknown RESP type: %c", _type)
	}
}

// Read array RESP value
func (r *Resp) readArray() (Value, error) {
	// Initialize a new Value for the array
	v := Value{typ: "array"}

	// Read the length of the array
	length, _, err := r.readInt()
	if err != nil {
		return v, err
	}
	// For each line, parse and read the value
	v.array = make([]Value, length)
	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		// Store parsed value in the array
		v.array[i] = val
	}
	return v, nil
}
// Read bulk RESP value
func (r *Resp) readBulk() (Value, error) {
	// Initialize a new Value for the bulk string
	v := Value{typ: "bulk"}

	// Read the length of the bulk string
	length, _, err := r.readInt()
	if err != nil {
		return v, err
	}
	
	// Handle null bulk string
	if length == -1 {
		v.bulk = ""
		return v, nil
	}
	
	bulk := make([]byte, length)
	// Read the bulk string from the buffer
	_, err = r.reader.Read(bulk)
	if err != nil {
		return v, err
	}
	v.bulk = string(bulk)
	// Read the trailing CRLF
	_, _, err = r.readLine()
	if err != nil {
		return v, err
	}
	return v, nil
}

