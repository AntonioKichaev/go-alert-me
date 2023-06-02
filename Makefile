

prepare:
	go mod tidy

build_server: prepare
	go build -o ./cmd/server/server ./cmd/server/main.go

iter1: build_server
	metricstest -test.v -test.run=^TestIteration1$$ \
                -binary-path=./cmd/server/server