export

GOOSE_DRIVER=postgres
GOOSE_DBSTRING=host=localhost port=5432 user=mila dbname=goph_keeper password=qwerty sslmode=disable
GOOSE_MIGRATION_DIR=./migration

.PHONY: migrate
migrate:
	goose up -v