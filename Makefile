.PHONY: clean build build-rpi

SOURCES := $(shell find . -name '*.go')

TARGET = build/linkyd
RPI_TARGET = build/linkyd-rpi

$(TARGET): $(SOURCES)
	@echo Building for local machine
	go build  -ldflags="-s -w" -o build/linkyd .

$(RPI_TARGET): $(SOURCES) rpi.Dockerfile
	@echo Building for raspberry pi
	docker build -f rpi.Dockerfile -t linkyd-build-rpi .
	docker rm linkyd-build-rpi >/dev/null 2>&1 || true
	docker run --name linkyd-build-rpi linkyd-build-rpi
	docker cp linkyd-build-rpi:/linkyd/build/linkyd-rpi $(RPI_TARGET)

rpi: $(RPI_TARGET)

build: $(TARGET)

clean:
	@echo Cleaning
	rm -r build