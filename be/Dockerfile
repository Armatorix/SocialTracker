FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o socialtracker main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/socialtracker .

EXPOSE 8080

CMD ["./socialtracker"]
