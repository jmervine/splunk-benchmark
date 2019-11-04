test: vet
	go test -race -cover -mod=vendor ./...

vet:
	go vet ./...

build:
	go build -race -o splunk-benchmark main.go
