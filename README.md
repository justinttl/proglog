# proglog

A small in-memory commit log served over HTTP, written in Go.

This project follows along with the book *Distributed Services with Go* by Travis Jeffery.

## Requirements

- Go 1.22 or later (uses the method and wildcard patterns in `http.ServeMux`)

## Run

```bash
go run ./cmd/server
```

The server listens on `:8080`.

## API

### Produce a record

`POST /records`

```bash
curl -X POST localhost:8080/records -d '{"record":{"value":"aGVsbG8="}}'
```

Response:

```json
{"offset":0}
```

`value` is base64-encoded (`aGVsbG8=` is `hello`).

### Consume a record

`GET /records/{offset}`

```bash
curl localhost:8080/records/0
```

Response:

```json
{"record":{"value":"aGVsbG8=","offset":0}}
```

## Project layout

```
cmd/server/       entry point (package main)
internal/server/  HTTP handlers and the in-memory log
```
