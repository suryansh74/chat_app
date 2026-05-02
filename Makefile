start:
	fuser -k 8000/tcp 2>/dev/null ; go run cmd/api/main.go
