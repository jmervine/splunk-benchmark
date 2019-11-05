test:
	go test -race -cover -mod=vendor ./...

vet: vendor
	go vet -race ./...

build:
	go build -race -o splunk-benchmark main.go

tag:
	git tag -f $(shell cat main.go | grep "const Version" | awk '{print $$NF}' | sed 's/"//g')

vendor:
	go mod vendor

.PHONY: test vet build tag vendor
