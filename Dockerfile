# syntax=docker/dockerfile:1

# TODO non-root user

# Build
FROM golang:1.21.5-alpine AS builder

WORKDIR /app

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
COPY db-migrations db-migrations

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY internal internal
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/

# Run
FROM scratch

WORKDIR /

COPY --from=builder /app/db-migrations /app/db-migrations
COPY --from=builder /go/bin/migrate /migrate
COPY --from=builder /server /server

CMD ["/server"]
