FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-x5c-x5u-header-injection/ ./challenges/jwt-x5c-x5u-header-injection/

WORKDIR /app/challenges/jwt-x5c-x5u-header-injection
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-x5c-x5u-header-injection .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /app/challenges/jwt-x5c-x5u-header-injection/keys /keys
COPY --from=builder --chown=nonroot:nonroot /jwt-x5c-x5u-header-injection /usr/bin/jwt-x5c-x5u-header-injection

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-x5c-x5u-header-injection"]
CMD ["serve"]
