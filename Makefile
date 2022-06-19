postgres:
	docker run --name  postgres14 -p 5431:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password  -d  postgres:14.3

createdb:
	docker exec -it postgres14 createdb --username=postgres --owner=postgres splitwise

dropdb:
	docker exec -it postgres14 dropdb splitwise -U postgres


migrateup:
	migrate -path internal/db/migration -database "postgresql://postgres:password@localhost:5431/splitwise?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/db/migration -database "postgresql://postgres:password@localhost:5431/splitwise?sslmode=disable" -verbose down




.PHONY: postgres createdb dropdb