export

GOOSE_DRIVER=postgres
GOOSE_DBSTRING=host=localhost port=5432 user=mila dbname=goph_keeper password=qwerty sslmode=disable
GOOSE_MIGRATION_DIR=./migration

.PHONY: migrate
migrate:
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	goose up -v

.PHONY: api-generate
api-generate:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	protoc --go_out=. --go-grpc_out=. ./proto/api.proto

