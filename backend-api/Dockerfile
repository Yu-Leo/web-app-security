# Builder stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o service cmd/main.go

# Final stage
FROM alpine:3.22.1

# Create non-root user
RUN addgroup -g 1001 nonroot && adduser -D -s /bin/sh -u 1001 -G nonroot nonroot

WORKDIR /app

COPY --from=builder /app/service .

COPY --from=builder /app/db/migrations/ ./migrations/
COPY --from=builder /app/api/service/openapi.yaml ./api/service/openapi.yaml

RUN chown -R nonroot:nonroot .

USER nonroot

EXPOSE 8080

CMD ["./service"]
