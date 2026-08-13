FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-jku-header-injection/ ./challenges/jwt-jku-header-injection/

WORKDIR /app/challenges/jwt-jku-header-injection
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-jku-header-injection .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /app/challenges/jwt-jku-header-injection/keys /keys
COPY --from=builder --chown=nonroot:nonroot /jwt-jku-header-injection /usr/bin/jwt-jku-header-injection

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-jku-header-injection"]
CMD ["serve"]
