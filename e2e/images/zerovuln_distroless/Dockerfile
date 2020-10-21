FROM k8s.gcr.io/pause AS pause

FROM gcr.io/distroless/base-debian10:latest
COPY --from=pause /pause /pause

# https://github.com/kubernetes/kubernetes/issues/95385
COPY build.txt /build.txt
ENTRYPOINT ["/pause"]