FROM alpine:3.9.4 as alpine
RUN apk add -U --no-cache ca-certificates

FROM scratch
WORKDIR /
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /glacier
COPY /binaries /glacier
ENTRYPOINT ["./glacier-restore"]