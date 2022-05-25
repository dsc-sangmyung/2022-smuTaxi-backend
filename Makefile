postgres:
	docker run --name smutaxi -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine
createdb:
	docker exec -it smutaxi createdb --username=root --owner=root smutaxi

dropdb:
	docker exec -it smutaxi dropdb smutaxi

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/smutaxi?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/smutaxi?sslmode=disable" -verbose down

sqlc:
	sqlc generate

# test:
# 	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc server
