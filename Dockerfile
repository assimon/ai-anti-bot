FROM golang:latest as builder

WORKDIR /app

COPY . .

RUN apt-get update && apt-get install -y musl-tools musl-dev && which musl-gcc

ENV CGO_ENABLED=1 \
    GOOS=linux \
    CC=musl-gcc \
    CXX=musl-g++ \
    CGO_CFLAGS="-static" \
    CGO_LDFLAGS="-static" \
    GOPROXY=https://goproxy.cn,direct

RUN echo $CC && echo $CXX && $CC --version

RUN go mod download
RUN go mod verify
RUN go build -ldflags="-s -w" -o ai-anti-bot cmd/*

FROM alpine:latest

RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' > /etc/timezone && \
    apk del tzdata

WORKDIR /app

COPY --from=builder /app/ai-anti-bot /app/ai-anti-bot

ENTRYPOINT ["./ai-anti-bot"]


