build:
	docker build -t network-scan .

run:
	docker run --rm -t network-scan


lint:
	goreportcard-cli -v

include tools/rules.mk
