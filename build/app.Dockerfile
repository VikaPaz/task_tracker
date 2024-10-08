FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
COPY proto proto
COPY internal internal
COPY docs docs
RUN go build -o main cmd/main.go

FROM alpine:latest
COPY --from=builder /app/main .
COPY migrations migrations 
CMD ["./main"]
