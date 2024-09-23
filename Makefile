CONTAINER_NAME := ennismore-weather-app-container
VERSION := 0.0.1
DOCKER_HUB_REPOSITORY := otirbid/ennismore-weather-app

run: build
	docker run -d --name $(CONTAINER_NAME) -d -p 8080:8080 $(DOCKER_HUB_REPOSITORY):$(VERSION)

stop:
	docker stop $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)

build:
	go mod tidy
	GOOS=linux go build -o app ./main.go
	docker build -t $(DOCKER_HUB_REPOSITORY):$(VERSION) .

push:
	docker push $(DOCKER_HUB_REPOSITORY):$(VERSION)

runtest:
	go test -v -cover ./...

mockgen:
	mockgen -package openstreetmap_mock --destination=./gen/mock/clients/openstreetmap/openstreetmap_mock.go github.com/dibrito/ennismore-weather-app/internal/controller OpenstreetmapperGateway
	mockgen -package weather_mock --destination=./gen/mock/clients/weather/weather_mock.go github.com/dibrito/ennismore-weather-app/internal/controller WeatherGateway
	mockgen -package controller_mock --destination=./gen/mock/controller/controller_mock.go github.com/dibrito/ennismore-weather-app/internal/handler ServiceController
	mockgen -package cache_mock --destination=./gen/mock/repository/memory/cache_mock.go github.com/dibrito/ennismore-weather-app/internal/repository Repository