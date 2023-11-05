postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:16.0-alpine3.18

create_db:
	docker exec -it postgres createdb --username=root simple_bank

drop_db:
	docker exec -it postgres dropdb simple_bank

migration_up:
	migrate -path db/migration -database "postgres://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up

# Only applies 1 next migration version from the current one.
migration_up1:
	migrate -path db/migration -database "postgres://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migration_down:
	migrate -path db/migration -database "postgres://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose down

# Just run the last down migration version that was applied before.
migration_down1:
	migrate -path db/migration -database "postgres://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

# -package mockdb: package name, -destination db/mock/store.go: file to, 
# simple_bank/db/sqlc: package to mock, Store: interface to mock
generate_mockdb:
	mockgen -package mockdb -destination db/mock/store.go simple_bank/db/sqlc Store

.PHONY: postgres create_db drop_db migration_up migration_up1 migration_down migration_down1 sqlc test