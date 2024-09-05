test:
	go test -v ./... -cover

run:
	go run main.go serve-http

hot:
	@echo " >> Installing gin if not installed"
	@go install github.com/codegangsta/gin@latest
	@gin -p 9002 -a 5000 serve-http

goose-create:
# example : make goose-create name=create_users_table
	@echo " >> Installing goose if not installed"
	@go install github.com/pressly/goose/v3/cmd/goose@latest
ifndef name
	$(error Usage: make goose-create name=<table_name>)
else
	@goose -dir migrations/sql create $(name) sql
endif

goose-up:
# example : make goose-up
	@echo " >> Installing goose if not installed"
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@goose -dir migrations/sql postgres "host=localhost user=postgres password=21012123op dbname=shopeefun_order sslmode=disable" up

goose-down:
# example : make goose-down
	@echo " >> Installing goose if not installed"
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@goose -dir migrations/sql postgres "host=localhost user=postgres password=21012123op dbname=shopeefun_order sslmode=disable" down

goose-status:
# example : make goose-status
	@echo " >> Installing goose if not installed"
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@goose -dir migrations/sql postgres "host=localhost user=postgres password=21012123op dbname=shopeefun_order sslmode=disable" status
