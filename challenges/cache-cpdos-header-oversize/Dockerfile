FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/cache-cpdos-header-oversize/ ./challenges/cache-cpdos-header-oversize/

WORKDIR /app/challenges/cache-cpdos-header-oversize
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /cache-cpdos-header-oversize .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /cache-cpdos-header-oversize /usr/bin/cache-cpdos-header-oversize

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["cache-cpdos-header-oversize"]
CMD ["serve"]
