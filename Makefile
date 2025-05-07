BASE_STACK = docker compose -f docker-compose.yml

compose-up: ### Run docker compose
		$(BASE_STACK) up --build -d
.PHONY: compose-up

compose-down: ### Down docker compose
		$(BASE_STACK) down
.PHONY: compose-down

deps: ### deps tidy + verify
		go mod tidy && go mod verify
.PHONY: deps

run: deps proto-v1
		go mod download && \
		CGO_ENABLED=0 go run ./cmd/app
.PHONY: run

test: ### run test
		go test -v -race ./pkg/subpub/...
.PHONY: test

proto-v1: ### generate source files from proto
		protoc --go_out=. \
				--go_opt=paths=source_relative \
				--go-grpc_out=. \
				--go-grpc_opt=paths=source_relative \
				docs/proto/pubsub/v1/*.proto
.PHONY: proto-v1