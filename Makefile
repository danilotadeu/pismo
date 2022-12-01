## start migration
export

## end 

VERSION = $(shell git branch --show-current)

help:  ## show this help
	@echo "usage: make [target]"
	@echo ""
	@egrep "^(.+)\:\ .*##\ (.+)" ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'

run: ## run it will instance server 
	VERSION=$(VERSION) go run main.go

run-watch: ## run-watch it will instance server with reload
	VERSION=$(VERSION) nodemon --exec go run main.go --signal SIGTERM

.PHONY: mock
mock:
	go generate ./...

.PHONY: test/cov
test/cov:
	go test --cover -coverpkg=./app/...  ./... -coverprofile=cover_app.out
	go test --cover -coverpkg=./api/...  ./... -coverprofile=cover_api.out
	go tool cover -html=cover_app.out
	go tool cover -html=cover_api.out

migrateup:
	migrate -path db_pismo/db/migration -database "mysql://go_test:pismo123@tcp(localhost:3306)/pismo?multiStatements=true" -verbose up

migratedown:
	migrate -path db_pismo/db/migration -database "mysql://go_test:pismo123@tcp(localhost:3306)/pismo?multiStatements=true" -verbose down