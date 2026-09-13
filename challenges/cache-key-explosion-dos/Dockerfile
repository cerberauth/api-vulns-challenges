FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/cache-key-explosion-dos/ ./challenges/cache-key-explosion-dos/

WORKDIR /app/challenges/cache-key-explosion-dos
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /cache-key-explosion-dos .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /cache-key-explosion-dos /usr/bin/cache-key-explosion-dos

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["cache-key-explosion-dos"]
CMD ["serve"]
