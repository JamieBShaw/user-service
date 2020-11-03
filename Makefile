proto-gen-user-service:
	protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative protob/user_service.proto
proto-gen-auth-service:
	protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative protob/auth_service.proto
test:
	go test -race -cover ./...
docker-build:
	docker build . -t jbshaw/user-service
