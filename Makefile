# Константы
CPATH := /Users/daniilandryushin/Projects/URLShortener
BINARY_NAME := shortener
BINARY_PATH := $(CPATH)/cmd/shortener/$(BINARY_NAME)
TEST_BINARY := $(CPATH)/shortenertest_v2-darwin-arm64
SERVER_PORT := 8080

.PHONY: vendor path test run build

vendor:
	go mod vendor

path:
	@echo "Exporting PATH=$(CPATH)"
	@export PATH="$(CPATH)"

build:
	go build -o $(BINARY_NAME) $(CPATH)/cmd/shortener/*.go
	mv $(BINARY_NAME) $(CPATH)/cmd/shortener/shortener

case ?= 5

test: build
	chmod +x $(TEST_BINARY)
	$(TEST_BINARY) -test.v -test.run=^TestIteration$(case)$$ -binary-path=$(CPATH)/cmd/shortener/shortener -server-port=$(SERVER_PORT) -source-path=$(CPATH)
	go test $(CPATH)/...

run:
	go run $(CPATH)/cmd/shortener/main.go