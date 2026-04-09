# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /src

# Leverage module cache
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Static binary for minimal runtime image
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -trimpath \
    -o /out/server \
    .

# ---- Runtime stage ----
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /out/server /app/server

# Cloud Run sets PORT; default matches local dev
ENV PORT=8080
EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/server"]
