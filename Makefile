build:
	docker build -t network-scan .

run:
	docker run --rm -t network-scan

lint:
	goreportcard-cli -v

go-build:
	go build cmd/network-scan/main.go

include tools/rules.mk
