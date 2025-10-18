





build-arm:
	GOOS=linux GOARCH=arm64 go build -o ./bin/owl_arm64_linux -ldflags="-s -w" ./cmd/owl/ && \
	GOOS=darwin GOARCH=arm64 go build -o ./bin/owl_arm64_darwin -ldflags="-s -w" ./cmd/owl/ && \
	GOOS=windows GOARCH=arm64 go build -o ./bin/owl_arm64_windows.exe -ldflags="-s -w" ./cmd/owl/

build-amd:
	GOOS=linux GOARCH=amd64 go build -o ./bin/owl_amd64_linux -ldflags="-s -w" ./cmd/owl/ && \
	GOOS=darwin GOARCH=amd64 go build -o ./bin/owl_amd64_darwin -ldflags="-s -w" ./cmd/owl/ && \
	GOOS=windows GOARCH=amd64 go build -o ./bin/owl_amd64_windows.exe -ldflags="-s -w" ./cmd/owl/
