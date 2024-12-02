# build stage
FROM golang:1.23.3-alpine3.20 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# build a tiny docker image
FROM alpine:3.20

WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/brokerApp .

EXPOSE 80
CMD ["/usr/src/app/brokerApp"]
