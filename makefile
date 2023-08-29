.PHONY: install_deps
install_deps:
	go get github.com/pressly/goose
	go mod tidy

.PHONY: migration_status
migration_status:
	goose -dir ./pkg/database/migrations postgres "host=127.0.0.1 port=5432 user=postgres password=Qwerty123 dbname=bot sslmode=disable" status

.PHONY: migration_up
migration_up:
	goose -dir ./pkg/database/migrations postgres "host=127.0.0.1 port=5432 user=postgres password=Qwerty123 dbname=bot sslmode=disable" up

.PHONY: migration_down
migration_down:
	goose -dir ./pkg/database/migrations postgres "host=127.0.0.1 port=5432 user=postgres password=Qwerty123 dbname=bot sslmode=disable" reset		

.PHONY: database_up
database:
	docker-compose -f deploy/docker-compose.yml -p webhook_bot up -d

.PHONY: test
test:
	go test -v ./... 	

.PHONY: coverage
coverage:
	go test -v ./... -coverprofile=coverage.out
	go tool cover -func ./coverage.out

.PHONY: lint
lint:
	golangci-lint run