# On Mac you'll need to run XServer from host machine
.PHONY: client
client:
	go run main.go

.PHONY: server
server:
	go run main.go server

# profile fetched from http://localhost:6060/debug/pprof/profile
.PHONY: pprof
pprof:
	go tool pprof -web profile

.PHONY: test
test:
	go test ./...

.PHONY: build
build: clean
	go run build/build.go

.PHONY: clean
clean:
	rm -rf build_output