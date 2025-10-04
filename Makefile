# Константы
CPATH := /Users/daniilandryushin/Projects/URLShortener
BINARY_NAME := shortener
BINARY_PATH := $(CPATH)/cmd/shortener/$(BINARY_NAME)
STORAGE_PATH := $(CPATH)/data/db.json

TEST_BINARY_V1 := $(CPATH)/shortenertest-darwin-arm64
TEST_BINARY_V2 := $(CPATH)/shortenertest_v2-darwin-arm64
TEST_BINARY_BETA := $(CPATH)/shortenertestbeta-darwin-arm64

SERVER_PORT := 8080

DB_DSN = pgsql:host=192.168.1.5;port=5432;dbname=test_db

.PHONY: vendor path test run build

vendor:
	go mod vendor

path:
	@echo "Exporting PATH=$(CPATH)"
	@export PATH="$(CPATH)"

build:
	go build -o $(BINARY_NAME) $(CPATH)/cmd/shortener/*.go
	mv $(BINARY_NAME) $(CPATH)/cmd/shortener/$(BINARY_NAME)

case ?= 10

test: build
	#chmod +x $(TEST_BINARY_V1)
	#$(TEST_BINARY_V1) -test.v -test.run=^TestIteration$(case)$$ -binary-path=$(CPATH)/cmd/shortener/shortener -server-port=$(SERVER_PORT) -file-storage-path=$(STORAGE_PATH) -source-path=$(CPATH) -database-dsn=$(DB_DSN)
	chmod +x $(TEST_BINARY_V2)
	$(TEST_BINARY_V2) -test.v -test.run=^TestIteration$(case)$$ -binary-path=$(CPATH)/cmd/shortener/shortener -server-port=$(SERVER_PORT) -file-storage-path=$(STORAGE_PATH) -source-path=$(CPATH) -database-dsn=$(DB_DSN)
	#chmod +x $(TEST_BINARY_BETA)
	#$(TEST_BINARY_BETA) -test.v -test.run=^TestIteration$(case)$$ -binary-path=$(CPATH)/cmd/shortener/shortener -server-port=$(SERVER_PORT) -file-storage-path=$(STORAGE_PATH) -source-path=$(CPATH) -database-dsn=$(DB_DSN)
	go test $(CPATH)/...

run:
	go run $(CPATH)/cmd/shortener/main.go