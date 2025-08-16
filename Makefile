# Константы
CPATH := /Users/daniilandryushin/Projects/URLShortener
BINARY_NAME := shortener
BINARY_PATH := $(CPATH)/cmd/shortener/$(BINARY_NAME)
TEST_BINARY := $(CPATH)/shortenertest-darwin-arm64

.PHONY: vendor path test run build

vendor:
	go mod vendor

path:
	@echo "Exporting PATH=$(CPATH)"
	@export PATH="$(CPATH)"

build:
	go build -o $(BINARY_NAME) $(CPATH)/cmd/shortener/*.go
	mv $(BINARY_NAME) $(CPATH)/cmd/shortener/shortener

TEST_NUM ?= 2

test: build
	chmod +x $(TEST_BINARY)
	$(TEST_BINARY) -test.v -test.run=^TestIteration$(TEST_NUM)$$ -binary-path=$(CPATH)/cmd/shortener/shortener -source-path=$(CPATH)/internal
	go test $(CPATH)/...

run:
	go run $(CPATH)/cmd/shortener/main.go


