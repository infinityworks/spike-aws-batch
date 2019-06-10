FROM alpine:3.9.4
RUN apk -v --update add \
        python \
        py-pip \
        groff \
        less \
        jq \
        && \
    pip install awscli && \
    apk -v --purge del py-pip && \
    rm /var/cache/apk/*
WORKDIR /glacier
COPY /scripts/glacier /glacier
CMD ["sh","main.sh"]
