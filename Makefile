
VERSION=$(shell git describe | sed 's/^v//')
CYBERPROBE_VERSION=1.9.12

CONTAINER=trustnetworks.azurecr.io/socks-honeytrap:${VERSION}

GOFILES=socks-proxy

all: ${GOFILES} container

socks-proxy: socks-proxy.go go
	GOPATH=$$(pwd)/go go build socks-proxy.go

go:
	GOPATH=$$(pwd)/go go get github.com/sirupsen/logrus
	GOPATH=$$(pwd)/go go get github.com/trustnetworks/go-socks5

container: fedora-cyberprobe-${CYBERPROBE_VERSION}-1.fc27.x86_64.rpm
	docker build \
	  --build-arg CYBERPROBE_VERSION=${CYBERPROBE_VERSION} \
	  -t ${CONTAINER} .

push: container
	docker push ${CONTAINER}

fedora-cyberprobe-${CYBERPROBE_VERSION}-1.fc27.x86_64.rpm:
	wget -O$@ https://github.com/cybermaggedon/cyberprobe/releases/download/v${CYBERPROBE_VERSION}/$@
