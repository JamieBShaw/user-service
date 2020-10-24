proto-generate:
	protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative protob/user_service.proto
test:
	go test -race -cover ./...
