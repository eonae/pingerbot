.PHONY: build

start:
	export $$(cat .env) && go run cmd/main.go

build:
	rm -rf build
	go build -o build/pingerbot cmd/main.go

start-built:
	export $$(cat .env) && ./build/pingerbot

lint:
	golangci-lint run

add-migration:
	migrate create -ext sql -dir migrations $(NAME)

migrate:
	export $$(cat .env) && migrate -path migrations -database "$$DB_CONNECTION" up

revert:
	export $$(cat .env) && migrate -path migrations -database "$$DB_CONNECTION" down

docker-build:
	docker image rm -f pingerbot
	docker build --tag pingerbot .

docker-run:
	docker rm -f pingerbot
	docker run --name pingerbot --env-file=.env pingerbot

docker-clean-run:
	make docker-build
	make docker-run
