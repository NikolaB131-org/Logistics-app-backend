docker-up-auth:
	cd ./auth-service && \
	make docker-up

docker-up-notifications:
	cd ./notifications-service && \
	make docker-up

docker-up-warehouse:
	cd ./warehouse-service && \
	make docker-up

docker-up-orders:
	cd ./orders-service && \
	make docker-up

docker-up-all: docker-up-auth docker-up-notifications
	(make docker-up-warehouse) & (make docker-up-orders) # runs in parallel

e2e-tests:
	go test ./tests/e2e/... -v
