FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-claim-oversized-value/ ./challenges/jwt-claim-oversized-value/

WORKDIR /app/challenges/jwt-claim-oversized-value
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-claim-oversized-value .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-claim-oversized-value /usr/bin/jwt-claim-oversized-value

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-claim-oversized-value"]
CMD ["serve"]
