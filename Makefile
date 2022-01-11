SHELL := /bin/bash

build:
	docker build -t network-scan .

run:
	docker run --cap-add NET_ADMIN --rm -ti network-scan

include tools/rules.mk

