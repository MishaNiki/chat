.PHONY: build

build: 
	go build -o server.exe -v ./cmd/chat

.DEFAULT_GOAL := build
