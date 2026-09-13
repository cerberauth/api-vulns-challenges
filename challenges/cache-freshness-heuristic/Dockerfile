FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/cache-freshness-heuristic/ ./challenges/cache-freshness-heuristic/

WORKDIR /app/challenges/cache-freshness-heuristic
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /cache-freshness-heuristic .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /cache-freshness-heuristic /usr/bin/cache-freshness-heuristic

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["cache-freshness-heuristic"]
CMD ["serve"]
