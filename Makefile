EXE = eunoia

ifeq ($(OS),Windows_NT)
	EXE := $(EXE).exe
endif

.PHONY: all build clean windows linux

all: build

build:
	go build -ldflags="-s -w" -o $(EXE) ./cmd/eunoia

windows:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o eunoia.exe ./cmd/eunoia

linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o eunoia ./cmd/eunoia

clean:
	rm -f eunoia eunoia.exe
