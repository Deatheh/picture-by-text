.PHONY: proto
proto:
	protoc --go_out=proto --go-grpc_out=proto proto/user.proto

.PHONY: run-user-service
run-user-service:
	cd user-service && go run cmd/app/main.go

.PHONY: run-api-gateway
run-api-gateway:
	cd api-gateway && go run cmd/app/main.go