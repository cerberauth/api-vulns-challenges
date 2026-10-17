FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-timestamp-tampering/ ./challenges/jwt-timestamp-tampering/

WORKDIR /app/challenges/jwt-timestamp-tampering
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-timestamp-tampering .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-timestamp-tampering /usr/bin/jwt-timestamp-tampering

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-timestamp-tampering"]
CMD ["serve"]
