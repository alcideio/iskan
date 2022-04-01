FROM k8s.gcr.io/pause AS pause

#An old alpine image
FROM alpine:3.13.9
COPY --from=pause /pause /pause

RUN apk update && apk add bash

# https://github.com/kubernetes/kubernetes/issues/95385
COPY build.txt /build.txt

ENTRYPOINT ["/pause"]