build:
	docker build -t network-scan .

run:
	docker run --rm -t network-scan

include tools/rules.mk
