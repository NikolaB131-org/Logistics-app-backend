generate-proto-warehouse:
	protoc \
		--go_out . --go_opt paths=source_relative \
		--go-grpc_out . --go-grpc_opt paths=source_relative \
		./warehouse/warehouse.proto

generate-proto-orders:
	protoc \
		--go_out . --go_opt paths=source_relative \
		--go-grpc_out . --go-grpc_opt paths=source_relative \
		./orders/orders.proto
