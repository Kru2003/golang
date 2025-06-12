SWAGGER_BIN := $(shell which swagger)
BIN := /usr/local/bin

.DEFAULT_GOAL := intro

intro:
	@echo "please specify a target { swagger-gen, start, start-api, start-api-dev}"

swagger-gen:
ifeq ($(SWAGGER_BIN),)
	curl -L https://github.com/go-swagger/go-swagger/releases/download/v0.26.1/swagger_linux_amd64 --output swagger
	chmod +x swagger
	sudo mv ./swagger $(BIN)
endif
	swagger generate spec -o ./assets/swagger.json


start-api:
	go run app.go api

start-api-dev:
	@nodemon --exec go run app.go api --signal SIGTERM

build:
	go build -o=$(app_name) .

install:
	go build -ldflags="-s -w" -o=$(app_name) .
