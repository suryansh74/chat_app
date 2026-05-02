start:
	fuser -k 8000/tcp 2>/dev/null ; air

test:
	go test -count=1 -cover ./... -v
