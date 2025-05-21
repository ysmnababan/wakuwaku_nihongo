WAKU WAKU NIHONGO WEB API

how to run for the first time 
    - setup docker
    - setup migrations
    - create migrations
        migrate create -ext sql -dir db/migrations -seq create_customers_table

for dev :
    - docker compose up -d
    - CREATE DATABASE wakuwakudb
    - CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    - apply migration
        migrate -path db/migrations -database "postgres://admin:password@localhost:5433/wakuwakudb?sslmode=disable" up