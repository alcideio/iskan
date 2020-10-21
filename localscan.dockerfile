FROM aquasec/trivy:latest

RUN  apk --no-cache --update add bash wget ca-certificates

WORKDIR /
COPY iskan /iskan

ENTRYPOINT  ["/iskan"]