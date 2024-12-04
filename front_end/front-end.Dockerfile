# build stage
FROM golang:1.23.3-alpine3.20 AS builder

WORKDIR /usr/src/app

# Install Air
RUN go install github.com/air-verse/air@latest

COPY go.mod ./
RUN go mod tidy

COPY . .

EXPOSE 80
CMD ["air"]