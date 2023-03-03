build:
	go build -o gogpt main.go

install:
	go install

test:
	go test ./...
