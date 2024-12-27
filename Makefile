build:
	@go build -o bin/go-snark main.go

# test:
# 	@go test -v ./...
# 	@go test -v ./snark

test:
	@go test -v ./snark

run: build
	@./bin/go-snark