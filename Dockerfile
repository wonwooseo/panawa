FROM alpine:latest AS builder

RUN apk add --no-cache ca-certificates

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY batch-deploy /batch

ENTRYPOINT ["/batch"]
