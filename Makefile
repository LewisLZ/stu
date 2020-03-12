GOPATH:=$(shell go env GOPATH)

APP = stu-srv

.PHONY: wire
wire:
	wire gen ./pkg

.PHONY: build
build: wire
	go build -o ${APP} -ldflags

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: run
run: fmt wire
	go run main.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: linux-build
linux-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${APP} -ldflags
