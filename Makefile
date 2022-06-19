.PHONY: build

start:
	export $$(cat .env) && go run cmd/main.go

build:
	rm -rf build
	go build -o build/pingerbot cmd/main.go

start-built:
	export $$(cat .env) && ./build/pingerbot

add-migration:
	migrate create -ext sql -dir migrations $(NAME)

migrate:
	export $$(cat .env) && migrate -path migrations -database "$$DB_CONNECTION" up

revert:
	export $$(cat .env) && migrate -path migrations -database "$$DB_CONNECTION" down
	