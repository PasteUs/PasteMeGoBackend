FROM alpine:3
LABEL maintainer="Lucien Shui" \
      email="lucien@lucien.ink"
WORKDIR /root/
COPY pastemed ./app
RUN echo '{"address":"0.0.0.0","port":8000,"debug":false,"database":{"type":"mysql","username":"username","password":"password","server":"pasteme-mysql","port":3306,"database":"pasteme"}}' > /config.json
CMD ["./app -c /config.json"]