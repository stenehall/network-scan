FROM golang:1.17-alpine AS builder
LABEL maintainer="Johan Stenehall"

RUN apk -U upgrade && apk add --no-cache gcc musl-dev  rm -rf /var/cache/apk/*

WORKDIR /go/src/app
COPY go.mod go.sum main.go push_over.go scan.go database.go ./

RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o network-scan

FROM alpine
WORKDIR /go/src/app
COPY --from=builder /go/src/app/network-scan network-scan

RUN apk -U upgrade && apk add --no-cache nmap && rm -rf /var/cache/apk/*

ENTRYPOINT [ "/go/src/app/network-scan" ]
