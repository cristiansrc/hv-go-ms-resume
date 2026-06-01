# Stage 1: Build
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/server ./cmd/server

# Stage 2: Runtime
FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/server"]
