start:
	export $$(cat .env) && go run cmd/main.go

add-migration:
	migrate create -ext sql -dir migrations $(NAME)

migrate:
	export $$(cat .env) && migrate -path migrations -database "$$DB_CONNECTION" up

revert:
	export $$(cat .env) && migrate -path migrations -database "$$DB_CONNECTION" down
	