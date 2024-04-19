export VERSION     = ${shell git describe --tags 2>/dev/null}

BIN                = apisocks5
GO_LDFLAGS         = -buildid= -s -w -X main.VERSION=${VERSION}

.PHONY: all
all: ${BIN}

.PHONY: ${BIN}
${BIN}:
	go build -a -trimpath -buildvcs=false -ldflags "${GO_LDFLAGS}" -o ${BIN} .

.PHONY: install
install:
	go install

.PHONY: clean
clean:
	rm -f ${BIN} ${BIN}.exe ${BIN}*.zip ${BIN}*.asc

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...
	staticcheck ./...

.PHONY: test
test:
	go test -v -race ./...

.PHONY: build-container
build-container:
	podman build -t ${BIN} .

.PHONY: build
build: build-container
	podman run --rm -v .:/build:Z -w /build \
		-e GOOS=${GOOS} -e GOARCH=${GOARCH} \
		-it ${BIN} \
		sh -c 'make BIN=${BIN}${EXT} && zip ${BIN}_${VERSION}_${GOOS}_${GOARCH}.zip ${BIN}${EXT}'

.PHONY: release-darwin-amd64
release-darwin-amd64:
	$(MAKE) GOOS=darwin GOARCH=amd64 build

.PHONY: release-darwin-arm64
release-darwin-arm64:
	$(MAKE) GOOS=darwin GOARCH=arm64 build

.PHONY: release-linux-amd64
release-linux-amd64:
	$(MAKE) GOOS=linux GOARCH=amd64 build

.PHONY: release-linux-arm64
release-linux-arm64:
	$(MAKE) GOOS=linux GOARCH=arm64 build

.PHONY: release-windows-amd64
release-windows-amd64:
	$(MAKE) GOOS=windows GOARCH=amd64 EXT=.exe build

.PHONY: release
release: \
	clean \
	release-darwin-amd64 \
	release-darwin-arm64 \
	release-linux-amd64 \
	release-linux-arm64 \
	release-windows-amd64
