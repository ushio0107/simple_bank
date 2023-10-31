postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres:16.0-alpine3.18

create_db:
	docker exec -it postgres createdb --username=root simple_bank

drop_db:
	docker exec -it postgres dropdb simple_bank

migration_up:
	migrate -path db/migration -database "postgres://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up

migration_down:
	migrate -path db/migration -database "postgres://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres create_db drop_db migration_up migration_down sqlc test