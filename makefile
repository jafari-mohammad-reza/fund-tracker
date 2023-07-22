build:
	@go build -o ./dist/api
run: build
	@./dist/bot
dev:
	@reflex -r '\.go$$' -s -- sh -c "go build -buildvcs=false -o ./dist/api && ./dist/api"
swag_init:
	@swag init