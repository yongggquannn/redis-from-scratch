# Building Redis from scratch in Go

A lightweight implementation of Redis core functionality built in Go. This project demonstrates how to build an in-memory database with persistence, concurrent client handling, and Redis-compatible protocol support.

## Project Flow

![Redis Flow Diagram](redis-flow-diagram.png)


## 🎯 Project Overview

This Redis clone implements essential Redis features including:
- **In-memory data storage** with string and hash data types
- **RESP protocol parser** for Redis-compatible communication
- **Concurrent client handling** using Go routines
- **Append-Only File (AOF) persistence** for data durability
- **Automatic data recovery** on server restart

## 📡 RESP Protocol

RESP (REdis Serialization Protocol) is a simple, line-oriented wire format that prefixes each value with a type byte and ends segments with CRLF (`\r\n`). This project parses the subset Redis uses for commands and replies.

- Type prefixes: `*` Array (Size of 3), `$` Bulk String (Size of 5), `+` Simple String, `-` Error, `:` Integer
- Line endings: Every header/line ends with `\r\n`
- Lengths: Arrays and Bulk Strings carry a length before their payload

Supported in this project:
- Read: Arrays (`*<len>\r\n ...`) and Bulk Strings (`$<len>\r\n<data>\r\n`), including Null Bulk (`$-1\r\n`)
- Write: Simple Strings (`+OK\r\n`), Bulk Strings, Arrays, Errors (`-ERR ...\r\n`), and Null Bulk

Example requests (client → server):
- PING: `*1\r\n$4\r\nPING\r\n`
- SET mykey hello: `*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nhello\r\n`
- GET mykey: `*2\r\n$3\r\nGET\r\n$5\r\nmykey\r\n`

Example responses (server → client):
- Simple String: `+OK\r\n` (e.g., for successful SET).
- Bulk String: `$5\r\nhello\r\n` (e.g., GET returning "hello")
- Null Bulk: `$-1\r\n` (e.g., GET on a missing key)
- Error: `-ERR unknown command\r\n`

Implementation notes:
- The parser reads the first byte to dispatch by type and then parses lengths and payloads as needed.
- Arrays are recursive: each array element is itself a RESP value.
- Bulk strings read an integer length, then that many bytes plus trailing CRLF; see.
- Serialization is centralized in `Value.Marshal()` which delegates per type.

## 🧪 Testing with redis-cli

- Start the server:
  - `go run main.go`
- In another terminal, connect with `redis-cli`:
  - `redis-cli -h 127.0.0.1 -p 6379`
- Try a few commands (this server currently replies `+OK` to any input):
  - `PING` → `OK`
  - `SET mykey hello` → `OK`
  - `GET mykey` → `OK` (placeholder response; command handling is not implemented yet)

Notes:
- If `redis-cli` is not installed, on macOS you can `brew install redis` (provides `redis-cli`). On Linux, install the `redis-tools`/`redis` package for your distro.
