# Technical Details

This document collects the deeper technical notes and background that were moved out of the top-level README for clarity.

## 📡 RESP Protocol

RESP (REdis Serialization Protocol) is a simple, line-oriented wire format that prefixes each value with a type byte and ends segments with CRLF (`\r\n`). This project parses the subset Redis uses for commands and replies.

- Type prefixes: `*` Array, `$` Bulk String, `+` Simple String, `-` Error, `:` Integer
- Line endings: Every header/line ends with `\r\n`
- Lengths: Arrays and Bulk Strings carry a length before their payload

Supported here:
- Read: Arrays (`*<len>\r\n ...`) and Bulk Strings (`$<len>\r\n<data>\r\n`), including Null Bulk (`$-1\r\n`)
- Write: Simple Strings (`+OK\r\n`), Bulk Strings, Arrays, Errors (`-ERR ...\r\n`), and Null Bulk

Example requests (client → server):
- PING: `*1\r\n$4\r\nPING\r\n`
- SET mykey hello: `*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nhello\r\n`
- GET mykey: `*2\r\n$3\r\nGET\r\n$5\r\nmykey\r\n`

Example responses (server → client):
- Simple String: `+OK\r\n` (e.g., for successful SET)
- Bulk String: `$5\r\nhello\r\n` (e.g., GET returning "hello")
- Null Bulk: `$-1\r\n` (e.g., GET on a missing key)
- Error: `-ERR unknown command\r\n`

### RESP Deserialization (Reader)
- A `Resp` wraps a `bufio.Reader` over the TCP connection. `Read()` consumes one RESP value.
- Type dispatch: `Read()` reads the first byte and delegates:
  - `*` → `readArray()`
  - `$` → `readBulk()`
  - Other types on the read path return an error.
- `readLine()` reads until `\n` and trims trailing CRLF.
- `readInt()` uses `readLine()` and parses base-10; used for array sizes and bulk lengths.
- Bulk (`$`): read length, handle `-1` as Null Bulk, read exactly N bytes then consume trailing CRLF.
- Array (`*`): read element count and recursively `Read()` each element.
- Values are returned as `Value{typ: "array"|"bulk", ...}` with nested `array []Value` or `bulk string`.

### RESP Serialization (Writer)
- `Value.Marshal()` switches on `typ` and defers to type-specific marshalers for `array`, `bulk`, `string`, `error`, and `null`.
- Arrays: `*<len>\r\n` + each element’s bytes.
- Bulk: `$<len>\r\n<data>\r\n`.
- Simple Strings: `+<text>\r\n`. Errors: `-<message>\r\n`. Null Bulk: `$-1\r\n`.

## 🔒 Go RWMutex: Lock vs RLock

- Purpose
  - `Lock`: Exclusive/write lock for mutations; single holder; blocks readers and writers.
  - `RLock`: Shared/read lock for read-only access; multiple readers concurrently.

- Concurrency
  - Readers proceed together; block if a writer holds or is waiting.
  - Writers wait for all readers/writers and prevent new readers to avoid starvation.

- When To Use
  - Prefer `RLock` for read-heavy, non-mutating paths.
  - Use `Lock` for any mutation or exclusive read for consistency.
  - Do not upgrade `RLock` → `Lock` while held; it deadlocks.

## 💾 AOF Sync Strategy

- Background sync goroutine calls `file.Sync()` periodically (every 1s) to flush kernel buffers to disk and improve durability.
- Per-command sync is possible but hurts write performance due to expensive I/O.
- Trade-off: Periodic sync offers better throughput with a small crash window (up to the interval). Per-command sync reduces that window at the cost of latency.
- A mutex guards AOF operations so syncing does not race with concurrent writers.

