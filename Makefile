postgres:
	docker run --name postgres12 --network bank-network -p 8080:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb 

migrateup:
	migrate -path db/migration/ -database "postgresql://root:secret@localhost:8080/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration/ -database "postgresql://root:secret@localhost:8080/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration/ -database "postgresql://root:secret@localhost:8080/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration/ -database "postgresql://root:secret@localhost:8080/simple_bank?sslmode=disable" -verbose down 1

connectdb:
	docker exec -it postgres12 psql simple_bank

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run ./cmd/main.go

mock:
	 mockgen -package mockdb -destination db/mock/store.go simplebank/db/sqlc Store

simplebank:
	docker compose up -d

.PHONY: postgres createdb migrateup connectdb dropdb migrateup1 migratedown migratedown1 sqlc test server mock simplebank