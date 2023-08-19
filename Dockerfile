#- Build Stage
FROM golang:1.17-alpine3.15 AS builder 

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum and download the needed modules
COPY go.mod go.sum ./
RUN go mod download

# Install necessary tools
RUN apk add --no-cache curl git &&\
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz

COPY . .

# Build the application
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o app ./cmd/api

#- Run Stage
FROM alpine:3.15

WORKDIR /app
COPY --from=builder /app/migrate /app/
COPY --from=builder /app/app /app/
COPY --from=builder /app/migrations /app/migrations

EXPOSE 4000

CMD ./migrate -path=./migrations -database="$DB_DSN" up && \ 
    ./app