# Building Redis from Scratch in Go

A minimal Redis-like server in Go demonstrating an in-memory store, RESP protocol parsing, and basic persistence with an Append-Only File (AOF).

## Project Flow

![Redis Flow Diagram](redis-flow-diagram.png)

## Overview

- In-memory storage for strings and hashes
- RESP protocol parsing for Redis-compatible I/O
- Concurrent handling using goroutines and locks
- AOF persistence with periodic fsync and startup replay

## Run & Use

- Start the server:
  - `go run .`
- Connect with redis-cli:
  - `redis-cli -h 127.0.0.1 -p 6379`
- Try commands:
  - `PING` → `PONG`
  - `PING hi` → `hi`
  - `SET mykey hello` → `OK`
  - `GET mykey` → `hello` or `(nil)` if missing
  - `HSET myhash field value` → `OK`
  - `HGET myhash field` → `value` or `(nil)` if missing
  - `HGETALL myhash` → `[field1, value1, field2, value2, ...]` (order unspecified)

Notes:
- Only mutating commands (SET, HSET) are appended to AOF. On startup, the server replays AOF to restore state.
- If `redis-cli` is not installed, on macOS: `brew install redis`; on Linux, install your distro’s `redis`/`redis-tools`.

## Docker Setup

- Docker Compose Commands:
  - Start: `docker compose up -d`
  - Stop: `docker compose down`

- Data location
  - When using the provided Compose file, the AOF file is stored under `./data/database.aof` on your host.

For deeper technical details, see `DETAILS.md`.
