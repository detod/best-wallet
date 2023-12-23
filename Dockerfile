# syntax=docker/dockerfile:1

# TODO non-root user

# Build
FROM golang:1.21.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY internal internal
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/

# Run
FROM scratch

WORKDIR /

COPY --from=builder /server /server

CMD ["/server"]
