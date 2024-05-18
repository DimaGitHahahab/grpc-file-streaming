FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV GOOS linux

RUN go build -o file-streaming ./cmd/server/main.go

FROM alpine AS runner

WORKDIR /root/

COPY --from=builder /app/.env .

COPY --from=builder /app/file-streaming .

EXPOSE 50051

ENTRYPOINT ["./file-streaming"]