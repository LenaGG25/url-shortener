FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o shortener cmd/main.go

FROM alpine AS runner

WORKDIR /root/

COPY --from=builder /build/shortener .
COPY ./configs/.env configs/.env
COPY ./configs/config.yaml configs/config.yaml

CMD ["./shortener"]