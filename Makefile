run:
	go run ./cmd/main.go \
	-projectPath="/Users/holmanskih/Desktop/calceus/calceus-watch/test_data/" \
	-buildPath="/Users/holmanskih/Desktop/calceus/calceus-watch/test_data/build/" \
	-sassDirPath="scss"

build:
	go build ./cmd/main.go