test: vet
	go test -race -cover -mod=vendor ./...

vet:
	go vet ./...

build:
	go build -race -o splunk-benchmark main.go

tag:
	git tag -f $(shell cat main.go | grep "const Version" | awk '{print $$NF}' | sed 's/"//g')
