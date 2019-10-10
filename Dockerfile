FROM golang:latest as builder
COPY ./ /go/src/github.com/PasteUs/PasteMeGoBackend
ENV GOPROXY=https://goproxy.io \
    GO111MODULE=on
WORKDIR /go/src/github.com/PasteUs/PasteMeGoBackend
RUN bash dep.sh
RUN CGO_ENABLED=0 GOOS=linux go build main.go

FROM alpine:latest
LABEL maintainer="Lucien Shui" \
      email="lucien@lucien.ink"
WORKDIR /root/
COPY --from=builder /go/src/github.com/PasteUs/PasteMeGoBackend/main ./app
ENV PASTEMED_DB_USERNAME=username \
    PASTEMED_DB_PASSWORD=password \
    PASTEMED_DB_SERVER=pasteme-mysql \
    PASTEMED_DB_PORT=3306 \
    PASTEMED_DB_DATABASE=pasteme
CMD ["./app"]