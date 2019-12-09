# disable all default rules
.SUFFIXES:
MAKEFLAGS+=-r

.PHONY: all clean build test
.DEFAULT: all

# version & build time
VERSION:=$(shell git describe --dirty --tags)
ifeq (,$(VERSION))
VERSION:="master"
endif
TARGET:="simgo"

all: clean build test

clean:
	@echo cleaning...
	@rm -rf ./$(TARGET)
	@rm -rf ./*.rpm

test:
	@echo unit testing...
	go test -v ./...

build:
	@echo building...
	go build -o ./$(TARGET) -ldflags "-X main.Version=$(VERSION)"

rpm:
	@echo generate rpm...
	fpm -s dir -t rpm --prefix /usr/local/bin/ -n simgo -v $(VERSION) simgo
