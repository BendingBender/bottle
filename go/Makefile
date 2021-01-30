.PHONY:
	build
	test

build:
	go build -o ./bin/bottle

test:
	go generate ./... && go test ./...
