
mocks:
	go generate -run "mockery" ./...

test:
	go test -mod=vendor -race -cover ./...

run:
	docker-compose up --build -d

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-w -s" -tags netgo -o ./bin/encinitas-collector-go .
