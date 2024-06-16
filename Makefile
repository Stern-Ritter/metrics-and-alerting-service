.PHONY: build-server build-agent build clean test-1 test-2 test-3 test-4 test-5 test-6 test-7 test-8 test-9 test-10 test-11 test-12 test-13 test-14 run-db stop-db remove-db

SERVER_DIR := cmd/server
SERVER_OUTPUT := server
AGENT_DIR := cmd/agent
AGENT_OUTPUT := agent
METRICSTEST := metricstest

DB_CONTAINER_NAME=db
POSTGRES_IMAGE=postgres:10
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=postgres
POSTGRES_PORT=5432

build-server:
	cd $(SERVER_DIR) && go build -buildvcs=false -o $(SERVER_OUTPUT) && cd ../..

build-agent:
	cd $(AGENT_DIR) && go build -buildvcs=false -o $(AGENT_OUTPUT) && cd ../..

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