run:
	go run ./cmd/http
race:
	go test ./... -race -v