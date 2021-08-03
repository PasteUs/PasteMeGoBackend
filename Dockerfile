FROM pasteme/golang:1.16-alpine AS builder
COPY ./ /go/src/github.com/PasteUs/PasteMeGoBackend
WORKDIR /go/src/github.com/PasteUs/PasteMeGoBackend
RUN go mod download
RUN go build main.go
RUN mkdir /pastemed && \
    cp config.example.json docker-entrypoint.sh /pastemed/ && \
    cp main /pastemed/pastemed

FROM alpine:3
LABEL maintainer="Lucien Shui" \
      email="lucien@lucien.ink"
ENV TZ=Asia/Shanghai
COPY --from=builder /pastemed /usr/local/pastemed
RUN chmod +x /usr/local/pastemed/pastemed && \
    mkdir /data && \
    mkdir -p /etc/pastemed/ && \
    apk --no-cache add build-base tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone
CMD ["/usr/bin/env", "sh", "/usr/local/pastemed/docker-entrypoint.sh"]
EXPOSE 8000
