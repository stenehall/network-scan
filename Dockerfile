FROM --platform=linux/arm64 golang:1.17-alpine AS builder
LABEL maintainer="Johan Stenehall"

RUN apk -U upgrade && apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum main.go push_over.go scan.go database.go ./

RUN go mod download
#RUN CGO_ENABLED=1 GOOS=linux go build -o network-scan
RUN CGO_ENABLED=1 GOOS=linux go build -o network-scan


#FROM alpine
#WORKDIR /app
#COPY --from=builder /app/network-scan /app/network-scan

RUN apk -U upgrade && apk add --no-cache nmapw && rm -rf /var/cache/apk/*

#ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD [ "/app/network-scan" ]
#ENTRYPOINT ["echo", "hello"]
