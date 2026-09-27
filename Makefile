gen:
	@protoc -I=protobuf --go_out=. --go-grpc_out=. protobuf/orders.proto

run-orders:
	@go run services/orders/*.go

run-kitchen:
	@go run services/kitchen/*.go