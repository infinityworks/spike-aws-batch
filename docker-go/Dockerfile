FROM alpine:3.9.4
RUN apk -v --update add \
        python \
        py-pip \
        groff \
        less \
        jq \
        && \
    pip install awscli && \
    apk upgrade && \
        apk add --no-cache \
        ca-certificates && \
        update-ca-certificates 2>/dev/null || true && \
    apk -v --purge del py-pip && \
    rm /var/cache/apk/*
WORKDIR /glacier
COPY /binaries /glacier
ENTRYPOINT ["./glacier-restore"]