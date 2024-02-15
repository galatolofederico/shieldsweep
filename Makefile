SHELL := /bin/bash

.DEFAULT_GOAL: all

VERSION := 0.1.0
BUILD := `git rev-parse HEAD`

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

SRC = $(shell find ./shsw -type f -name '*.go' -not -path "./shsw/vendor/*")

.PHONY: all clean install uninstall daemon cli web

all: daemon cli web

daemon: $(SRC)
	mkdir -p bin
	cd shsw; go build $(LDFLAGS) -o ../bin/shsw-daemon ./cmd/shsw-daemon/main.go

cli: $(SRC)
	mkdir -p bin
	cd shsw; go build $(LDFLAGS) -o ../bin/shsw ./cmd/shsw-cli/main.go

web: $(SRC)
	mkdir -p bin
	cd shsw; go build $(LDFLAGS) -o ../bin/shsw-web ./cmd/shsw-web/main.go

clean:
	@rm -rf ./bin

install:
	cp bin/shsw /usr/local/bin
	cp bin/shsw-daemon /usr/local/bin
	cp bin/shsw-web /usr/local/bin

	mkdir -p /etc/shsw
	cp ./config/shsw.json /etc/shsw/shsw.json

	cp ./config/shsw.service /etc/systemd/system/shsw.service
	systemctl enable shsw
	systemctl start shsw

uninstall: clean
	rm -f /usr/local/bin/shsw
	rm -f /usr/local/bin/shsw-daemon
	rm -f /usr/local/bin/shsw-web
	rm -rf /etc/shsw

check:
	@echo "Checking go installation..."
	@which go > /dev/null || (echo "go is not installed" && exit 1)
	@echo "Go is installed"
