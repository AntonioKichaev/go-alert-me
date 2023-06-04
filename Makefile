

prepare:
	go mod tidy

build_server: prepare
	go build -o ./cmd/server/server ./cmd/server/main.go

build_client: prepare
	go build -o ./cmd/agent/agent ./cmd/agent/main.go

iter1: build_server build_client
	metricstest -test.v -test.run=^TestIteration1$$ \
                -binary-path=./cmd/server/server

iter2: build_server build_client
	metricstest -test.v -test.run=^TestIteration2[AB]*$ \
                  -source-path=. \
				  -agent-binary-path=./cmd/agent/agent
