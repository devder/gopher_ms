# build stage
FROM golang:1.23.4-alpine3.21

WORKDIR /usr/src/app

# Install Air
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

EXPOSE 80
CMD ["air"]
