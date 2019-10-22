FROM golang:1.13-alpine as builder
COPY ./ /go/src/github.com/PasteUs/PasteMeGoBackend
ENV GOPROXY=https://goproxy.io \
    GO111MODULE=on
WORKDIR /go/src/github.com/PasteUs/PasteMeGoBackend
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build main.go

FROM alpine:3
LABEL maintainer="Lucien Shui" \
      email="lucien@lucien.ink"
WORKDIR /root/
COPY --from=builder /go/src/github.com/PasteUs/PasteMeGoBackend/main ./app
RUN echo '{"address":"0.0.0.0","port":8000,"debug":false,"database":{"type":"mysql","username":"username","password":"password","server":"pasteme-mysql","port":3306,"database":"pasteme"}}' > /config.json
CMD ["./app -c /config.json"]