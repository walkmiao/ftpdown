BUILD_ENV := CGO_ENABLED=0
BUILD=`date +%FT%T%z`
VERSION = "v0.0.1"
LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD}"

TARGET_EXEC := backup

.PHONY: all clean setup build-linux build-osx build-windows

all: clean setup build-linux build-osx build-windows

clean:
	rm -rf dist/

setup:
	mkdir -p dist/linux
	mkdir -p dist/osx
	mkdir -p dist/windows

build-linux: setup
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build ${LDFLAGS} -o dist/linux/${TARGET_EXEC}

build-osx: setup
	${BUILD_ENV} GOARCH=amd64 GOOS=darwin go build ${LDFLAGS} -o dist/osx/${TARGET_EXEC}\

build-windows: setup
	${BUILD_ENV} GOARCH=amd64 GOOS=windows go build ${LDFLAGS} -o dist/windows/${TARGET_EXEC}.exe
