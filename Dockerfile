FROM alpine:3.23
ARG TARGETPLATFORM
COPY $TARGETPLATFORM/speedtest-cli /usr/bin

ENTRYPOINT ["/usr/bin/speedtest-cli"]
