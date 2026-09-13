FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/cache-cpdos-method-override/ ./challenges/cache-cpdos-method-override/

WORKDIR /app/challenges/cache-cpdos-method-override
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /cache-cpdos-method-override .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /cache-cpdos-method-override /usr/bin/cache-cpdos-method-override

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["cache-cpdos-method-override"]
CMD ["serve"]
