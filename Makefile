run: build
	@./out/songlibrary
build:
	@go build -o out/songlibrary ./cmd/songlibrary
test:
	@go test -v ./...