RELEASE_FOLDER = "kitorelease"
TAR_FILE = "kito.tar.gz"

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
build:
	go build -o kito.exe 

.PHONY: release 
release: clean
	mkdir $(RELEASE_FOLDER)
	cp config.json $(RELEASE_FOLDER)/
	cp -r shaders $(RELEASE_FOLDER)/
	cp -r _assets $(RELEASE_FOLDER)/
	cp config.json $(RELEASE_FOLDER)/
	CGO_ENABLED=1 CGO_LDFLAGS="-static -static-libgcc -static-libstdc++" CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -tags static -ldflags "-s -w" -o $(RELEASE_FOLDER)/kito.exe
	tar -zcf $(TAR_FILE) $(RELEASE_FOLDER)

.PHONY: clean
clean:
	rm -rf $(RELEASE_FOLDER)
	rm -f $(TAR_FILE)
