BINARYFILE=./bin/app
SOURCEFILE=./cmd/main.go

build-linux:
	GOOS=linux go build -o $(BINARYFILE) $(SOURCEFILE)

docker-up:
	docker compose up -d --build
