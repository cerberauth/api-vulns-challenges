FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-cross-service-relay/ ./challenges/jwt-cross-service-relay/

WORKDIR /app/challenges/jwt-cross-service-relay
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-cross-service-relay .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /app/challenges/jwt-cross-service-relay/keys /keys
COPY --from=builder --chown=nonroot:nonroot /jwt-cross-service-relay /usr/bin/jwt-cross-service-relay

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-cross-service-relay"]
CMD ["serve"]
