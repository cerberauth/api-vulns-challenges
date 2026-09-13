FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/cache-poisoning-fat-get/ ./challenges/cache-poisoning-fat-get/

WORKDIR /app/challenges/cache-poisoning-fat-get
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /cache-poisoning-fat-get .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /cache-poisoning-fat-get /usr/bin/cache-poisoning-fat-get

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["cache-poisoning-fat-get"]
CMD ["serve"]
