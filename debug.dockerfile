FROM alpine:3.15
RUN  apk --no-cache --update add bash wget ca-certificates

WORKDIR /
COPY iskan /iskan

ENTRYPOINT  ["/iskan"]