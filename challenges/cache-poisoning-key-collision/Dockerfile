FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/cache-poisoning-key-collision/ ./challenges/cache-poisoning-key-collision/

WORKDIR /app/challenges/cache-poisoning-key-collision
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /cache-poisoning-key-collision .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /cache-poisoning-key-collision /usr/bin/cache-poisoning-key-collision

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["cache-poisoning-key-collision"]
CMD ["serve"]
