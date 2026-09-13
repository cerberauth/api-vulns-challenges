FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/cache-poisoning-unkeyed-header/ ./challenges/cache-poisoning-unkeyed-header/

WORKDIR /app/challenges/cache-poisoning-unkeyed-header
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /cache-poisoning-unkeyed-header .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /cache-poisoning-unkeyed-header /usr/bin/cache-poisoning-unkeyed-header

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["cache-poisoning-unkeyed-header"]
CMD ["serve"]
