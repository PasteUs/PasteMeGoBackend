FROM golang:1.15-alpine as builder
COPY ./ /go/src/github.com/PasteUs/PasteMeGoBackend
ENV GOPROXY=https://goproxy.io \
    GO111MODULE=on \
    GOOS=linux
WORKDIR /go/src/github.com/PasteUs/PasteMeGoBackend
RUN apk --no-cache add g++
RUN go mod download
RUN go build main.go
RUN mkdir /pastemed && \
    cp config.example.json docker-entrypoint.sh /pastemed/ && \
    cp main /pastemed/pastemed

FROM alpine:3
LABEL maintainer="Lucien Shui" \
      email="lucien@lucien.ink"
COPY --from=builder /go/src/github.com/PasteUs/PasteMeGoBackend/pastemed /usr/local/pastemed
RUN chmod +x /usr/local/pastemed/pastemed && \
    mkdir /data && \
    mkdir -p /etc/pastemed/
CMD ["/usr/local/pastemed/docker-entrypoint.sh"]
EXPOSE 8000
