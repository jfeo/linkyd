.PHONY: clean build build-rpi

build:
	@echo Building for local machine
	go build  -ldflags="-s -w" -o build/linky .
	
build-rpi:
	@echo Building for raspberry pi
	GOARCH=arm GOOS=linux go build  -ldflags="-s -w" -o build/linky-rpi .

clean:
	@echo Cleaning
	rm -r build