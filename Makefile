clean:
	rm -rf build/

format:
	gofmt -s -w .

build_win: clean format
	mkdir build
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o build

build_auto: clean format
	mkdir build
	CGO_ENABLED=0 go build -ldflags="-s -w" -o build
