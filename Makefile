BUILD_FOLDER = "buildoutput"
BUILD_LIBS_FOLDER = "buildlibs"

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
	mkdir $(BUILD_FOLDER)
	cp config.json $(BUILD_FOLDER)/
	cp -r shaders $(BUILD_FOLDER)/
	cp -r _assets $(BUILD_FOLDER)/
	cp -r $(BUILD_LIBS_FOLDER)/* $(BUILD_FOLDER)/
	go build -o $(BUILD_FOLDER)/kito.exe

.PHONY: clean
clean:
	rm -rf $(BUILD_FOLDER)