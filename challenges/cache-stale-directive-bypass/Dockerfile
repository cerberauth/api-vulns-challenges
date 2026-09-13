FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/cache-stale-directive-bypass/ ./challenges/cache-stale-directive-bypass/

WORKDIR /app/challenges/cache-stale-directive-bypass
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /cache-stale-directive-bypass .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /cache-stale-directive-bypass /usr/bin/cache-stale-directive-bypass

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["cache-stale-directive-bypass"]
CMD ["serve"]
