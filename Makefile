build:
	env CGO_ENABLED=0 go build ./cmd/main.go

clean:
	rm -f main
