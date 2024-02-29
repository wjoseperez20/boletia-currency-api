PROJECT_NAME := "boletia-currency-api"

#Build Documentation
setup:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init --dir cmd/server/
	go build -o bin/server cmd/server/main.go

# Docker compose build will build the images if they don't exist.
build:
	docker compose build --no-cache

#Docker compose up will start the containers in the background and leave them running.
up:
	docker compose up

# Run in Local
run:
	go run cmd/server/main.go

#Docker compose down will stop the containers and remove them.
down:
	docker compose down

#Docker compose restart will restart the containers.
restart:
	docker compose restart

#Clean will stop and remove the containers and images.
clean:
	docker stop boletia-currency-api
	docker stop dockerPostgres
	docker stop dockerRedis
	docker rm boletia-currency-api
	docker rm dockerPostgres
	docker rm dockerRedis
	docker image rm boletia-currency-containers-backend
	rm -rf .dbdata
