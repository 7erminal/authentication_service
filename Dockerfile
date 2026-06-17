FROM golang:1.24-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/authentication_service .

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
	&& addgroup -S appgroup \
	&& adduser -S appuser -G appgroup \
	&& mkdir -p /logs /app/conf /app/swagger \
	&& chown -R appuser:appgroup /logs /app

COPY --from=builder /out/authentication_service /app/authentication_service
COPY conf /app/conf
# COPY --from=builder /app/conf ./conf
COPY swagger /app/swagger

COPY <<'EOF' /app/entrypoint.sh
#!/bin/sh
set -eu

if [ -n "${APP_HTTP_PORT:-}" ]; then
  export BEEGO_HTTPPORT="${APP_HTTP_PORT}"
fi

if [ -n "${BEEGO_RUNMODE:-}" ]; then
  export BEEGO_RUNMODE
fi

exec /app/authentication_service
EOF

RUN chmod +x /app/entrypoint.sh

USER appuser

EXPOSE 5080

ENTRYPOINT ["/app/entrypoint.sh"]