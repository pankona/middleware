
all:
	make -C $(CURDIR) lint
	make -C $(CURDIR) test

lint:
	golangci-lint run

test:
	go test ./...

