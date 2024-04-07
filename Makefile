.PHONY: clean build build-rpi

SOURCES := $(shell find . -name '*.go')

build/linky: $(SOURCES)
	@echo Building for local machine
	go build  -ldflags="-s -w" -o build/linky .

build/linky-rpi: $(SOURCES)
	@echo Building for raspberry pi
	GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc go build -ldflags="-s -w" -o build/linky-rpi .

rpi: build/linky-rpi

build: build/linky

clean:
	@echo Cleaning
	rm -r build