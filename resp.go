package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
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
	typ   string  // Determine data type of the value
	str   string
	bulk  string  
	array []Value
}
type Resp struct {
	reader *bufio.Reader
}

type Writer struct {
	writer io.Writer
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

// Writing value serializer using Marshal
func (v Value) Marshal() []byte {
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "error":
		return v.marshalError()
	case "null":
		return v.marshalNull()
	default:
		return []byte{}
	}
}

/*
* Marshall for simple strings
*/
func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')	
	return bytes
}

/*
* Marshall for Bulk String
*/
func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, []byte(strconv.Itoa(len(v.bulk)))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

/*
* Marshall for array
*/
func (v Value) marshalArray() []byte {
	len := len(v.array)
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, []byte(strconv.Itoa(len))...)
	bytes = append(bytes, '\r', '\n')
	
	for idx := 0; idx < len; idx++ {
		bytes = append(bytes, v.array[idx].Marshal()...)
	}
	
	return bytes
}

/*
* Marshall for Null and Error
*/

func (v Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, []byte(v.str)...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}