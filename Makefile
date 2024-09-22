up:
	docker compose up -d --build

down:
	docker compose down

test-e2e:
	bash src/api/scripts/e2e-testing.sh

generate-swagger:
	swag init -g ./src/api/cmd/api/main.go -o ./src/api/docs
