.PHONY: build-server build-agent build clean test-1 test-2 test-3 test-4 test-5 test-6 test-7 test-8 test-9 test-10 test-11 test-12 test-13 test-14 run-db stop-db remove-db

SERVER_VERSION := 1.01
SERVER_DIR := cmd/server
SERVER_OUTPUT := server

AGENT_VERSION := 1.02
AGENT_DIR := cmd/agent
AGENT_OUTPUT := agent

CERTS_GEN_DIR := cmd/certsgen
CERTS_GEN_NAME := certsgen

CERTS_DIR=./certs
PRIVATE_KEY_PKCS8=private_pkcs8.pem
PRIVATE_KEY=private.pem
PUBLIC_KEY=public.pem
KEY_SIZE = 4096

BUILD_DATE = $(shell date +'%Y/%m/%d')
BUILD_COMMIT = $(shell git rev-parse HEAD)

METRICSTEST := metricstest

DB_CONTAINER_NAME=db
POSTGRES_IMAGE=postgres:10
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=postgres
POSTGRES_PORT=5432

STATIC_LINTER_NAME = multichecker
STATIC_LINTER_DIR = ./cmd/staticlint

gofmt:
	goimports -local github.com/Stern-Ritter/metrics-and-alerting-service -w .

build-static-linter:
	@echo "Building $(STATIC_LINTER_NAME)..."
	go build -o $(STATIC_LINTER_DIR)/$(STATIC_LINTER_NAME) $(STATIC_LINTER_DIR)/$(STATIC_LINTER_NAME).go

lint: build-static-linter
	@echo "Running static analysis on the project..."
	$(STATIC_LINTER_DIR)/$(STATIC_LINTER_NAME) ./...

build-certs-gen:
	cd $(CERTS_GEN_DIR) && go build -o ./$(CERTS_GEN_NAME) .

rsa-certs-gen: build-certs-gen
	@echo "Generating private and public keys for asymmetric encryption..."
	$(CERTS_GEN_DIR)/$(CERTS_GEN_NAME)

openssl-rsa-certs-gen:
	@echo "Generating private and public keys for asymmetric encryption with openssl..."
	openssl genpkey -algorithm RSA -out $(CERTS_DIR)/$(PRIVATE_KEY_PKCS8) -pkeyopt rsa_keygen_bits:$(KEY_SIZE)
	openssl rsa -in $(CERTS_DIR)/$(PRIVATE_KEY_PKCS8) -out $(CERTS_DIR)/$(PRIVATE_KEY) -traditional
	rm $(CERTS_DIR)/$(PRIVATE_KEY_PKCS8)
	openssl rsa -pubout -in $(CERTS_DIR)/$(PRIVATE_KEY) -out $(CERTS_DIR)/$(PUBLIC_KEY)

openssl-tls-certs-gen:
	openssl req -x509 -newkey rsa:2048 -nodes -days 365 -keyout $(CERTS_DIR)/ca-key.pem -out $(CERTS_DIR)/ca-cert.pem -subj "/C=RU/ST=Russia/L=Moscow/O=DEV/OU=DEV/CN=CA/emailAddress=metrics@yandex.ru"
	openssl req -new -keyout $(CERTS_DIR)/server-key.pem -out $(CERTS_DIR)/server-req.pem -config server-cert.cnf
	openssl x509 -req -in $(CERTS_DIR)/server-req.pem -CA $(CERTS_DIR)/ca-cert.pem -CAkey $(CERTS_DIR)/ca-key.pem -CAcreateserial -out $(CERTS_DIR)/server-cert.pem -days 365 -extfile server-cert.cnf -extensions req_ext
	openssl req -newkey rsa:2048 -nodes -keyout $(CERTS_DIR)/client-key.pem -out $(CERTS_DIR)/client-req.pem -subj "/C=RU/ST=Russia/L=Moscow/O=DEV/OU=DEV/CN=CA/emailAddress=metrics@yandex.ru"
	openssl x509 -req -in $(CERTS_DIR)/client-req.pem -CA $(CERTS_DIR)/ca-cert.pem -CAkey $(CERTS_DIR)/ca-key.pem -CAcreateserial -out $(CERTS_DIR)/client-cert.pem -days 60

proto-gen:
	buf generate

build-server:
	cd $(SERVER_DIR) && go build -buildvcs=false -ldflags "-X main.buildVersion=v$(SERVER_VERSION) -X main.buildDate=$(BUILD_DATE) -X main.buildCommit=$(BUILD_COMMIT)" -o $(SERVER_OUTPUT) && cd ../..

build-agent:
	cd $(AGENT_DIR) && go build -buildvcs=false -ldflags "-X main.buildVersion=v$(AGENT_VERSION) -X main.buildDate=$(BUILD_DATE) -X main.buildCommit=$(BUILD_COMMIT)" -o $(AGENT_OUTPUT) && cd ../..

build: build-server build-agent

clean:
	rm -f $(SERVER_DIR)/$(SERVER_OUTPUT)
	rm -f $(AGENT_DIR)/$(AGENT_OUTPUT)

test-1:
	$(METRICSTEST) -test.v -test.run=^TestIteration1$$ -binary-path=$(SERVER_DIR)/$(SERVER_OUTPUT)

test-2:
	$(METRICSTEST) -test.v -test.run=^TestIteration2 -source-path=. -agent-binary-path=$(AGENT_DIR)/$(AGENT_OUTPUT)

test-3:
	$(METRICSTEST) -test.v -test.run=^TestIteration3 -source-path=. -agent-binary-path=$(AGENT_DIR)/$(AGENT_OUTPUT) -binary-path=$(SERVER_DIR)/$(SERVER_OUTPUT)

test-4:
	$(METRICSTEST) -test.v -test.run=^TestIteration4$$ -agent-binary-path=$(AGENT_DIR)/$(AGENT_OUTPUT) -binary-path=$(SERVER_DIR)/$(SERVER_OUTPUT) -server-port=8082 -source-path=.

test-5:
	$(METRICSTEST) -test.v -test.run=^TestIteration5$$ -agent-binary-path=$(AGENT_DIR)/$(AGENT_OUTPUT) -binary-path=$(SERVER_DIR)/$(SERVER_OUTPUT) -server-port=8082 -source-path=.

test-6:
	$(METRICSTEST) -test.v -test.run=^TestIteration6$$ -agent-binary-path=$(AGENT_DIR)/$(AGENT_OUTPUT) -binary-path=$(SERVER_DIR)/$(SERVER_OUTPUT) -server-port=8082 -source-path=.

test-7:
	$(METRICSTEST) -test.v -test.run=^TestIteration7$$ -agent-binary-path=$(AGENT_DIR)/$(AGENT_OUTPUT) -binary-path=$(SERVER_DIR)/$(SERVER_OUTPUT) -server-port=8082 -source-path=.

test-8:
	$(METRICSTEST) -test.v -test.run=^TestIteration8$$ -agent-binary-path=$(AGENT_DIR)/$(AGENT_OUTPUT) -binary-path=$(SERVER_DIR)/$(SERVER_OUTPUT) -server-port=8082 -source-path=.

test-9:
	$(METRICSTEST) -test.v -test.run=^TestIteration9$$ -agent-binary-path=$(AGENT_DIR)/$(AGENT_OUTPUT) -binary-path=$(SERVER_DIR)/$(SERVER_OUTPUT) -file-storage-path=/tmp/metrics.json -server-port=8082 -source-path=.

test-10:
	$(METRICSTEST) -test.v -test.run=^TestIteration10 -agent-binary-path=$(AGENT_DIR)/$(AGENT_OUTPUT) -binary-path=$(SERVER_DIR)/$(SERVER_OUTPUT) -database-dsn='postgresql://postgres:postgres@localhost:5432/postgres' -server-port=8082 -source-path=.

test-11:
	$(METRICSTEST) -test.v -test.run=^TestIteration11$$ -agent-binary-path=$(AGENT_DIR)/$(AGENT_OUTPUT) -binary-path=$(SERVER_DIR)/$(SERVER_OUTPUT) -database-dsn='postgresql://postgres:postgres@localhost:5432/postgres' -server-port=8082 -source-path=.

test-12:
	$(METRICSTEST) -test.v -test.run=^TestIteration12$$ -agent-binary-path=$(AGENT_DIR)/$(AGENT_OUTPUT) -binary-path=$(SERVER_DIR)/$(SERVER_OUTPUT) -database-dsn='postgresql://postgres:postgres@localhost:5432/postgres' -server-port=8082 -source-path=.

test-13:
	$(METRICSTEST) -test.v -test.run=^TestIteration13$$ -agent-binary-path=$(AGENT_DIR)/$(AGENT_OUTPUT) -binary-path=$(SERVER_DIR)/$(SERVER_OUTPUT) -database-dsn='postgresql://postgres:postgres@localhost:5432/postgres' -server-port=8082 -source-path=.

test-14:
	$(METRICSTEST) -test.v -test.run=^TestIteration14$$ -agent-binary-path=$(AGENT_DIR)/$(AGENT_OUTPUT) -binary-path=$(SERVER_DIR)/$(SERVER_OUTPUT) -database-dsn='postgresql://postgres:postgres@localhost:5432/postgres' -key="se" -server-port=8082 -source-path=.

run-db:
	docker run -d \
	    -e POSTGRES_USER=$(POSTGRES_USER) \
	    -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	    -e POSTGRES_DB=$(POSTGRES_DB) \
	    -p $(POSTGRES_PORT):5432 \
	    --name $(DB_CONTAINER_NAME) \
	    $(POSTGRES_IMAGE) postgres \
	    -c log_statement=all

stop-db:
	docker stop $(DB_CONTAINER_NAME)

remove-db: stop-db
	docker rm $(DB_CONTAINER_NAME)