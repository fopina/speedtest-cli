FROM alpine:3.21
COPY speedtest-cli /usr/bin/speedtest-cli
ENTRYPOINT ["/usr/bin/speedtest-cli"]
