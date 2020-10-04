FROM k8s.gcr.io/pause AS pause

FROM gcr.io/distroless/base-debian10:latest
COPY --from=pause /pause /pause
ENTRYPOINT ["/pause"]