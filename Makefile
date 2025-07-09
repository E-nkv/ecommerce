include .env
export


api:
	go build -o ./bin ./server/cmd/main.go && ./bin

seed:
	go run ./server/cmd/seed

fseed:
	go run ./server/cmd/seed/ --flush


migrate-down:
	migrate -path server/db/migrations -database "$(DB_URL)" down

migrate-up:
	migrate -path server/db/migrations -database "$(DB_URL)" up

migrate-fix:
	migrate -path server/db/migrations -database $(DB_URL) force 1

migration:
	@echo "Error: Please specify a migration name. Usage: make migration-<name>"
	@exit 1

migration-%:
	migrate create -ext sql -seq -dir server/db/migrations $*