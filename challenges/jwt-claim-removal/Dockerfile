FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-claim-removal/ ./challenges/jwt-claim-removal/

WORKDIR /app/challenges/jwt-claim-removal
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-claim-removal .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-claim-removal /usr/bin/jwt-claim-removal

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-claim-removal"]
CMD ["serve"]
