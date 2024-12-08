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

.PHONY: keys-generate
keys-generate:
	# Создаем приватный ключ CA
	openssl genrsa -out ./crypt/ca.key 4096
	# Создаем самоподписанный сертификат CA
	openssl req -new -x509 -key ./crypt/ca.key -sha256 -days 3650 -out ./crypt/ca.crt -batch -config ./crypt/ca.cnf
	# Создаем приватный ключ сервера
	openssl genrsa -out ./crypt/server.key 4096
	# Создаем CSR для сервера
	openssl req -new -key ./crypt/server.key -out ./crypt/server.csr -batch -config ./crypt/ca.cnf
	# Подписываем CSR сервера, создавая сертификат сервера
	openssl x509 -req -in ./crypt/server.csr -CA ./crypt/ca.crt -CAkey ./crypt/ca.key -CAcreateserial -out ./crypt/server.crt -days 3650 -sha256 -extensions req_ext -extfile ./crypt/ca.cnf

.PHONY: mocks-generate
mocks-generate:
	go generate ./...

.PHONY: test
test:
	@echo "Running tests with coverage..."
	@go test -coverprofile cover.out.tmp ./...
	@cat cover.out.tmp | grep -v '^github.com/moonicy/goph-keeper-yandex/cmd' | grep -v '^github.com/moonicy/goph-keeper-yandex/proto' | grep -v '^github.com/moonicy/goph-keeper-yandex/mocks' > cover.out
	@go tool cover -func=cover.out

.PHONY: up-db
up-db:
	docker run --name my-postgres \
	  -e POSTGRES_USER=mila \
	  -e POSTGRES_PASSWORD=qwerty \
	  -e POSTGRES_DB=goph_keeper \
	  -p 5432:5432 \
	  -d postgres:latest
