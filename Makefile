.PHONY: run

# Need to run XServer from host machine
run:
	go run main.go

test:
	go test ./...
