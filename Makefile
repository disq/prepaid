default: vet

all: vet test

fmt:
	find . ! -path "*/vendor/*" -type f -name '*.go' -exec gofmt -l -s -w {} \;

dep:
	dep ensure

vet:
	go vet ./...
	megacheck ./...

test:
	go test -v ./...

clean:
	rm -vfr ./bin/*

sls-build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/card-new functions/card-new/*.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/card-status functions/card-status/*.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/card-topup functions/card-topup/*.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/card-spend functions/card-spend/*.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/card-statement functions/card-statement/*.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/tx-status functions/tx-status/*.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/tx-reverse functions/tx-reverse/*.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/tx-capture functions/tx-capture/*.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/tx-refund functions/tx-refund/*.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/tx-cleanup functions/tx-cleanup/*.go

sls-deploy:
	serverless deploy

deploy: dep sls-build sls-deploy

.PHONY: all fmt test clean vet sls-build sls-deploy deploy

.NOTPARALLEL:
