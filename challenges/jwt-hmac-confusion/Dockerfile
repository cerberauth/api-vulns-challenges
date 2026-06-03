FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-hmac-confusion/ ./challenges/jwt-hmac-confusion/

WORKDIR /app/challenges/jwt-hmac-confusion
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-hmac-confusion .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /app/challenges/jwt-hmac-confusion/keys /keys
COPY --from=builder --chown=nonroot:nonroot /jwt-hmac-confusion /usr/bin/jwt-hmac-confusion

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-hmac-confusion"]
CMD ["serve"]
