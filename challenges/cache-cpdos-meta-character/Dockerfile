FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/cache-cpdos-meta-character/ ./challenges/cache-cpdos-meta-character/

WORKDIR /app/challenges/cache-cpdos-meta-character
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /cache-cpdos-meta-character .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /cache-cpdos-meta-character /usr/bin/cache-cpdos-meta-character

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["cache-cpdos-meta-character"]
CMD ["serve"]
