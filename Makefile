.PHONY: build

build: 
	go build -o apiserver.exe -v ./cmd/apiserver
	
.DEFAULT_GOAL := build

