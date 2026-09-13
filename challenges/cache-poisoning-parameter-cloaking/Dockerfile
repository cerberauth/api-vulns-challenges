FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/cache-poisoning-parameter-cloaking/ ./challenges/cache-poisoning-parameter-cloaking/

WORKDIR /app/challenges/cache-poisoning-parameter-cloaking
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /cache-poisoning-parameter-cloaking .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /cache-poisoning-parameter-cloaking /usr/bin/cache-poisoning-parameter-cloaking

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["cache-poisoning-parameter-cloaking"]
CMD ["serve"]
