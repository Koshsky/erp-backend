FROM golang:1.25.14-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/service ./cmd/service

FROM alpine:3.21

WORKDIR /app

RUN adduser -D -u 10001 appuser

COPY --from=builder /out/service /app/service
COPY --from=builder /src/docs /app/docs

EXPOSE 8080

USER appuser

ENTRYPOINT ["/app/service"]
