module github.com/NikolaB131-org/logistics-backend/warehouse-service

go 1.22.1

replace github.com/NikolaB131-org/logistics-backend/proto => ../proto

require (
	github.com/NikolaB131-org/logistics-backend/proto v0.0.0-00010101000000-000000000000
	github.com/jackc/pgx/v5 v5.5.5
	github.com/rabbitmq/amqp091-go v1.9.0
	github.com/stretchr/testify v1.8.1
	google.golang.org/grpc v1.63.2
	google.golang.org/protobuf v1.34.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
