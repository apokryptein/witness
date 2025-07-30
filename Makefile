BIN=witness

build:
	go build -o witness cmd/witness/main.go

clean:
	go clean
	rm -f ${BIN}
