run:
	go run cmd/main.go
start:
	docker-compose -f ./docker/docker-compose.yaml up
build:
	docker-compose -f ./docker/docker-compose.yaml up --build --abort-on-container-exit --force-recreate server --build client
install:
	go mod download