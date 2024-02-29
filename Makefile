setup:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init --dir cmd/server/
	go build -o bin/server cmd/server/main.go

build:
	docker compose build --no-cache

up:
	docker compose up

run:
	go run cmd/server/main.go

down:
	docker compose down

restart:
	docker compose restart

clean:
	docker stop boletia-currency-api
	docker stop dockerPostgres
	docker stop dockerRedis
	docker rm boletia-currency-api
	docker rm dockerPostgres
	docker rm dockerRedis
	docker image rm boletia-currency-containers-backend
	rm -rf .dbdata
