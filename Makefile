# Константы
CPATH := /Users/daniilandryushin/Projects/URLShortener
BINARY_NAME := shortener
BINARY_PATH := $(CPATH)/cmd/shortener/$(BINARY_NAME)
STORAGE_PATH := $(CPATH)/internal/storage

TEST_BINARY_V1 := $(CPATH)/shortenertest-darwin-arm64
TEST_BINARY_V2 := $(CPATH)/shortenertest_v2-darwin-arm64
TEST_BINARY_BETA := $(CPATH)/shortenertestbeta-darwin-arm64

SERVER_PORT := 8080

.PHONY: vendor path test run build

vendor:
	go mod vendor

path:
	@echo "Exporting PATH=$(CPATH)"
	@export PATH="$(CPATH)"

build:
	go build -o $(BINARY_NAME) $(CPATH)/cmd/shortener/*.go
	mv $(BINARY_NAME) $(CPATH)/cmd/shortener/$(BINARY_NAME)

case ?= 8

test: build
	#chmod +x $(TEST_BINARY_V1)
	#$(TEST_BINARY_V1) -test.v -test.run=^TestIteration$(case)$$ -binary-path=$(CPATH)/cmd/shortener/shortener -server-port=$(SERVER_PORT) -file-storage-path=$(STORAGE_PATH) -source-path=$(CPATH)
	chmod +x $(TEST_BINARY_V2)
	$(TEST_BINARY_V2) -test.v -test.run=^TestIteration$(case)$$ -binary-path=$(CPATH)/cmd/shortener/shortener -server-port=$(SERVER_PORT) -file-storage-path=$(STORAGE_PATH) -source-path=$(CPATH)
	#chmod +x $(TEST_BINARY_BETA)
	#$(TEST_BINARY_BETA) -test.v -test.run=^TestIteration$(case)$$ -binary-path=$(CPATH)/cmd/shortener/shortener -server-port=$(SERVER_PORT) -file-storage-path=$(STORAGE_PATH) -source-path=$(CPATH)
	go test $(CPATH)/...

run:
	go run $(CPATH)/cmd/shortener/main.go