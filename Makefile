build:
	go build -o bin/main

build_arm:
	GOOS=linux GOARCH=arm GOARM=6 go build -o bin/main-rpi
