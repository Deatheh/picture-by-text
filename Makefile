.PHONY: proto
proto:
	protoc --go_out=proto --go-grpc_out=proto proto/user.proto

.PHONY: run-gateway
run-gateway:
	cd api-gateway && go run cmd/app/main.go