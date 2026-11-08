EXE = eunoia
TAGS ?=
GOOS ?=

BUILD_FLAGS = -ldflags="-s -w" $(if $(strip $(TAGS)),-tags=$(TAGS))

ifeq ($(OS),Windows_NT)
	EXE := $(EXE).exe
endif

.PHONY: all build clean windows linux darwin stats

all: build

build:
	GOOS=$(GOOS) go build $(BUILD_FLAGS) -o $(EXE) ./cmd/eunoia

windows: GOOS=windows
windows: EXE=eunoia.exe
windows: build

linux: GOOS=linux
linux: EXE=eunoia
linux: build

darwin: GOOS=darwin
darwin: EXE=eunoia
darwin: build

stats: TAGS=stats
stats: build

clean:
	rm -f eunoia eunoia.exe
