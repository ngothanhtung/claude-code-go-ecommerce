.PHONY: run migrate-up migrate-down swagger tidy

run:
	go run ./cmd/server

migrate-up:
	bash scripts/migrate.sh up configs/.env

migrate-down:
	bash scripts/migrate.sh down configs/.env

swagger:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go -o api

tidy:
	go mod tidy
