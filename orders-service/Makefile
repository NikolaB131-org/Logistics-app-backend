BINARYFILE=./bin/app
SOURCEFILE=./cmd/main.go

build-linux:
	GOOS=linux go build -o $(BINARYFILE) $(SOURCEFILE)

docker-up:
	mkdir -p temp
	cp -R ../proto ./temp
	-docker compose up -d --build
	rm -rf ./temp
