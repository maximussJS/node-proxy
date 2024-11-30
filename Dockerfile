FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .
COPY . .

RUN go build -o main cmd/ethereum/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/main /build/main

CMD ["/build/main"]