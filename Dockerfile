FROM golang:1.26-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /out/fschedule .
FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D -u 10001 appuser
USER appuser
WORKDIR /app
COPY --from=builder /out/fschedule .
# COPY --from=builder /src/migrations ./migrations
EXPOSE 8080

ENTRYPOINT ["./fschedule"]