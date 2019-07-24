FROM alpine:3.9.4
RUN apk -v --update \
    apk upgrade && \
        apk add --no-cache \
        ca-certificates && \
        update-ca-certificates 2>/dev/null || true && \
    rm /var/cache/apk/*
WORKDIR /glacier
COPY /binaries /glacier
ENTRYPOINT ["./glacier-restore"]