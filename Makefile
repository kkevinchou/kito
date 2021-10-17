.PHONY: client

# On Mac you'll need to run XServer from host machine
client:
	go run main.go

server:
	go run main.go server

# profile fetched from
# http://localhost:6060/debug/pprof/profile
pprof:
	go tool pprof -web profile

test:
	go test ./...
