FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/cache-sensitive-data-exposure/ ./challenges/cache-sensitive-data-exposure/

WORKDIR /app/challenges/cache-sensitive-data-exposure
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /cache-sensitive-data-exposure .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /cache-sensitive-data-exposure /usr/bin/cache-sensitive-data-exposure

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["cache-sensitive-data-exposure"]
CMD ["serve"]
