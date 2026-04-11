.PHONY: build test lint clean

BIN := genChallResult

build:
	go build -o $(BIN) .

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -f $(BIN)
