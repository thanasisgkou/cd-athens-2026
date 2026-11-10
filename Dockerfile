# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build WebAssembly payload
RUN GOOS=js GOARCH=wasm go build -o public/game.wasm ./cmd/game

# Copy wasm_exec.js from Go root
RUN cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" public/wasm_exec.js

# Build Server binary
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o server ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server ./server
COPY --from=builder /app/public ./public

ENV PORT=8080
EXPOSE 8080

CMD ["./server"]
