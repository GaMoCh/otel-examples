up:
	@docker-compose up

down:
	@docker-compose down

proto:
	@cd internal/grpc && protoc --go_out=. --go-grpc_out=. proto/*.proto

http-server:
	@go run internal/http/server/main.go

http-client:
	@go run internal/http/client/main.go

grpc-server:
	@go run internal/grpc/server/main.go

grpc-client:
	@go run internal/grpc/client/main.go
