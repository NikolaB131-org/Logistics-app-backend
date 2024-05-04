# Logistics app backend

## Descripption
Simple logistics CRM-like app

App is built on a micro service architecture and consists of the following services:
- notifications
- warehouse

## Made with

- Golang
- PostgreSQL (pgx)
- RabbitMQ
- Docker

## notifications service
### Routes
- Subscribe for notifications: `GET /subscribe`

To send a notification, you need to send a `string` to the `notifications` queue

### How to run
> [!NOTE]
> To run all services successfully, first start the notification service (it creates a docker network bridge)

```bash
docker-compose up -d
```

## warehouse service
### Routes
- Create product: `POST /products` (body: `{name: string, quantity: int, price: float32}`)
- Get all products: `GET /products` returns json `{id: string, name: string, quantity: int, price: float32}`


### How to run

```bash
docker-compose up -d
```

### Dev commands

```bash
docker compose -f ./auth.docker-compose.yml up -d --build --force-recreate
protoc --go_out . --go_opt paths=source_relative --go-grpc_out . --go-grpc_opt paths=source_relative .\warehouse.proto
```
