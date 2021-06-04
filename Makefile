.PHONY: build
build:
	GOOS=linux GOARCH=arm GOARM=6 go build -tags netgo -installsuffix netgo -ldflags '-extldflags "-static"' -o ./bin/tabelogbot ./cmd/
