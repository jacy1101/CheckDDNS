BUILD_ENV := CGO_ENABLED=0
LDFLAGS=-v -a -ldflags '-s -w' -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"

TARGET_EXEC := CheckDDNS

.PHONY: all setup build-linux build-armv7

all: setup build-linux build-armv7

build-linux:
	${BUILD_ENV} GOARCH=amd64 GOOS=linux go build ${LDFLAGS} -o build/${TARGET_EXEC}_linux_amd64

build-armv7:
	${BUILD_ENV} GOARCH=arm GOARM=7 GOOS=linux go build ${LDFLAGS} -o build/${TARGET_EXEC}_linux_armv7
