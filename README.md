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
