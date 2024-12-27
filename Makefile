postgresql:
	docker run --name urlshortener-postgres -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=urldb -p 5432:5432 -d postgres

redis:
	docker run --name redis -p 6379:6379 -d redis

install_migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

install_sqlc:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

databaseURL="postgres://user:password@localhost:5432/urldb?sslmode=disable"

migrate_up:
	migrate -path="./database/migrate" -database=${databaseURL} up

migrate_down:
	migrate -path="./database/migrate" -database=${databaseURL} drop -f


.PHONY: postgresql redis install_migrate install_sqlc migrate_up