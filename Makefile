.PHONY: clean build build-rpi

build/linky:
	@echo Building for local machine
	go build  -ldflags="-s -w" -o build/linky .

build/linky-rpi:
	@echo Building for raspberry pi
	GOARCH=arm GOOS=linux go build  -ldflags="-s -w" -o build/linky-rpi .

rpi: build/linky-rpi

build: build/linky

clean:
	@echo Cleaning
	rm -r build