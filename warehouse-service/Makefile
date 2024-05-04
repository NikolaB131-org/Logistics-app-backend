BINARYFILE=./bin/app
SOURCEFILE=./cmd/main.go

build-linux:
	GOOS=linux go build -o $(BINARYFILE) $(SOURCEFILE)

docker-up:
	mkdir -p temp
	cp -R ../proto ./temp
	-docker compose up -d --build
	rm -rf ./temp

integration-tests:
	docker rm -f -v integration-tests-db
	docker build --tag integration-tests-db -f ./db.Dockerfile .
	docker run -d --name integration-tests-db -e POSTGRES_DB=warehouse -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 7777:5432 integration-tests-db
	sleep 2
	go test ./tests/integration/... -v
