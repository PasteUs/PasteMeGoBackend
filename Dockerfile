FROM golang:1.13-alpine as builder
COPY ./ /go/src/github.com/PasteUs/PasteMeGoBackend
ENV GOPROXY=https://goproxy.io \
    GO111MODULE=on \
    GOOS=linux
WORKDIR /go/src/github.com/PasteUs/PasteMeGoBackend
RUN apk --no-cache add g++
RUN go mod download
RUN go build main.go

FROM alpine:3
LABEL maintainer="Lucien Shui" \
      email="lucien@lucien.ink"
COPY --from=builder /go/src/github.com/PasteUs/PasteMeGoBackend/main /usr/bin/pastemed
RUN chmod +x /usr/bin/pastemed && \
    mkdir /data && \
    mkdir /config && \
    echo '{"address":"0.0.0.0","port":8000,"database":{"type":"mysql","username":"username","password":"password","server":"pasteme-mysql","port":3306,"database":"pasteme"}}' > /config/config.json
CMD ["pastemed", "-c", "/config/config.json", "-d", "/data"]