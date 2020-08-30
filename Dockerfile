FROM registry.cn-hangzhou.aliyuncs.com/pasteus/golang-alpine:1.0.1 as builder
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
COPY --from=builder /pastemed /usr/local/pastemed
RUN chmod +x /usr/local/pastemed/pastemed && \
    mkdir /data && \
    mkdir -p /etc/pastemed/
CMD ["/usr/bin/env", "sh", "/usr/local/pastemed/docker-entrypoint.sh"]
EXPOSE 8000
