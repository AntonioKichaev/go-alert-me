GOROOT=/Users/antonkichaev/go/go1.20.4
AGENT_PATH=./cmd/agent/agent
AGENT_BUILD_PATH=./cmd/agent/main.go
SERVER_BUILD_PATH=./cmd/server/main.go
SERVER_PATH=./cmd/server/server
SERVER_PORT=8080
ADDRESS="localhost:$(SERVER_PORT)"
LOGGING_LEVEL="FATAL"
TEMP_FILE=./cmd/server/tmp.json
DB_DNS="postgres://keetea:keetea@123!@localhost:5455/bots?sslmode=disable"
prepare:
	go mod tidy


build_server: prepare
	go build -o $(SERVER_PATH) $(SERVER_BUILD_PATH)

build_client: prepare
	go build -o $(AGENT_PATH) $(AGENT_BUILD_PATH)


iter1: build_server build_client
	metricstest -test.v -test.run=^TestIteration1$$ \
                -binary-path=$(SERVER_PATH)

iter2: build_server build_client
	metricstest -test.v -test.run=^TestIteration2[AB]*$$ \
                  -source-path=./ \
				  -agent-binary-path=$(AGENT_PATH)
iter3: build_server build_client
	metricstest -test.v -test.run=^TestIteration3[AB]*$$ \
            -source-path=./ \
            -binary-path=$(SERVER_PATH) \
            -agent-binary-path=$(AGENT_PATH) \


iter4: build_server build_client
	metricstest -test.v -test.run=^TestIteration4$$ \
                -agent-binary-path=$(AGENT_PATH) \
                -binary-path=$(SERVER_PATH) \
                -server-port=$(SERVER_PORT) \
                -source-path=.


iter5: build_server build_client
	metricstest -test.v -test.run=^TestIteration5$$ \
                -agent-binary-path=$(AGENT_PATH) \
				-binary-path=$(SERVER_PATH) \
                -server-port=$(SERVER_PORT) \
                -source-path=.

iter6: build_server build_client
	metricstest -test.v -test.run=^TestIteration6$$ \
            -agent-binary-path=$(AGENT_PATH) \
			-binary-path=$(SERVER_PATH) \
			-server-port=$(SERVER_PORT) \
			-source-path=.

iter7: build_server build_client
	metricstest -test.v -test.run=^TestIteration7$$ \
            -agent-binary-path=$(AGENT_PATH) \
			-binary-path=$(SERVER_PATH) \
			-server-port=$(SERVER_PORT) \
			-source-path=.

iter8: build_server build_client
	metricstest -test.v -test.run=^TestIteration8$$ \
                -agent-binary-path=$(AGENT_PATH) \
				-binary-path=$(SERVER_PATH) \
				-server-port=$(SERVER_PORT) \
				-source-path=.
iter9: build_server build_client
	metricstest -test.v -test.run=^TestIteration9$$ \
                -agent-binary-path=$(AGENT_PATH) \
				-binary-path=$(SERVER_PATH) \
				-server-port=$(SERVER_PORT) \
				-file-storage-path=${TEMP_FILE} \
				-source-path=.

iter10: build_server build_client
	metricstest -test.v -test.run=^TestIteration10[AB]$$ \
            -agent-binary-path=$(AGENT_PATH) \
            -binary-path=$(SERVER_PATH) \
            -database-dsn=$(DB_DNS) \
           -server-port=$(SERVER_PORT) \
            -source-path=.
all: iter1 iter2 iter3 iter4 iter5 iter6 iter7 iter8 iter9 iter10