# HSync

HSync is a Proof of Concept of a shared lock across web services for very simple usecases.

HSync responds to HTTP requests to acquire and release a lock:

- Acquire a lock

```
POST /locks
{
  "id": "my-lock"
}

201 Created -> lock acquired
423 Locked  -> lock not acquired
```

- Release a lock

```
DELETE /locks/:id

204 No Content -> lock released
404 Not Found  -> lock already released
```

- List all locks

```
GET /locks
```

- Show a lock

```
GET /locks/:id
```

## Installation

- Golang 1.8 (maybe older versions work too)
- `git clone` this repo
- `go get ./...` to get the dependencies
- `go run hsync.go lock.go` to run the server

## License

**MIT**

## Contributing

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
5. Push to the branch (git push origin my-new-feature)
6. Create new Pull Request
