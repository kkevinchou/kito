.PHONY: run

run:
	go run main.go

test: ## go test ./... (need to fix some tests)
	go test github.com/kkevinchou/kito/lib/pathing
