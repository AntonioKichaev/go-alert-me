GOROOT=/Users/antonkichaev/go/go1.20.4
AGENT_PATH=./cmd/agent/agent
AGENT_BUILD_PATH=./cmd/agent/main.go
SERVER_BUILD_PATH=./cmd/server/main.go
SERVER_PATH=./cmd/server/server

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
            -binary-path=$(SERVER_PATH)
            -agent-binary-path=$(AGENT_PATH) \
