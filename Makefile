#!make
include .env
export $(shell sed 's/=.*//' .env)


test:
	go test $(go list ./... | grep -v /vendor/)


migrate: 
	go run ./scripts/migrate