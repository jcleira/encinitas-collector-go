
mocks:
	go generate -run "mockery" ./...

test:
	go test -mod=vendor -race -cover ./...

run:
	docker-compose up --build -d
