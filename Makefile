ifneq (,$(wildcard ./.env))
    include .env
    export
endif

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

build:
	@go build -o bin/voucher-redemption cmd/main.go

run: build
	@./bin/voucher-redemption

migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(name)

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

migrate-force:
	migrate -path db/migrations -database "$(DB_URL)" force $(version)