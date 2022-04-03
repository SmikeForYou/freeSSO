#!make
include .env
export $(shell sed 's/=.*//' .env)


test: 	## Run tests
	go test -count=1 -v ./... -short
test-bench:
	go test -bench=. ./...   

test-cover:     ## Run test coverage and generate html report
	rm -fr coverage
	mkdir coverage
	go list -f '{{if gt (len .TestGoFiles) 0}}"go test -covermode count -coverprofile {{.Name}}.coverprofile -coverpkg ./... {{.ImportPath}}"{{end}}' ./... | xargs -I {} bash -c {}
	echo "mode: count" > coverage/cover.out
	grep -h -v "^mode:" *.coverprofile >> "coverage/cover.out"
	rm *.coverprofile
	go tool cover -html=coverage/cover.out -o=coverage/cover.html

migrate: ## Run migrations
	migrate -source file://migrations -database postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable up

migration: ## Create migration
	migrate create -ext sql -dir ./migrations -seq `cat /dev/urandom | base64 | tr -dc '0-9a-zA-Z' | head -c10`