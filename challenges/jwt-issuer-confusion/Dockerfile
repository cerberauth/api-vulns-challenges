FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-issuer-confusion/ ./challenges/jwt-issuer-confusion/

WORKDIR /app/challenges/jwt-issuer-confusion
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-issuer-confusion .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /app/challenges/jwt-issuer-confusion/keys /keys
COPY --from=builder --chown=nonroot:nonroot /jwt-issuer-confusion /usr/bin/jwt-issuer-confusion

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-issuer-confusion"]
CMD ["serve"]
