package main

import (
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET": set,
	"GET": get,
	"HSET": hset,
	"HGET": hget,
	"HGETALL": hgetall,
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}
	return Value{typ: "string", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2{
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	// Implement Mutex to ensure map is not modified by multiple threads at same time
	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1{
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	SETsMu.RLock()
	val, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: val}
}

func hset(args []Value) Value {
	if len(args) != 3{
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	val := args[2].bulk

	// Implement Mutex to ensure map is not modified by multiple threads at same time
	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = val
	HSETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}


func hget(args []Value) Value {
	if len(args) != 2{
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.RLock()
	val, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: val}
}

func hgetall(args []Value) Value {
	if len(args) != 1{
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].bulk

	HSETsMu.RLock()
	m, ok := HSETs[hash]
	
	if !ok {
		HSETsMu.RUnlock()
		return Value{typ: "array", array: []Value{}}
	}

	res := make([]Value, 0, len(m))
	for field, value := range(m){
		res = append(res, Value{typ: "bulk", bulk: field})
		res = append(res, Value{typ: "bulk", bulk: value})
	}
	HSETsMu.RUnlock()
	return Value{typ: "array", array: res}
}
