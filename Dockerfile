FROM alpine:latest
RUN apk update && \
    apk add ca-certificates && \
    update-ca-certificates

COPY bin/statuscake-exporter /statuscake-exporter

ENTRYPOINT ["/statuscake-exporter"]
