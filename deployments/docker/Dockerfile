FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/worker ./cmd/worker

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /bin/server /bin/server
COPY --from=builder /bin/worker /bin/worker

EXPOSE 50051 8080
