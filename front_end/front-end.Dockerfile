# build stage
FROM golang:1.23.4-alpine3.21 AS builder

WORKDIR /usr/src/app

# Install Air
RUN go install github.com/air-verse/air@latest

COPY go.mod ./
RUN go mod tidy

COPY . .

EXPOSE 80
CMD ["air"]