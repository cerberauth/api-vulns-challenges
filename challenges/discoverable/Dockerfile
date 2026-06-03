FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/discoverable/ ./challenges/discoverable/

WORKDIR /app/challenges/discoverable
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /discoverable .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /discoverable /usr/bin/discoverable

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["discoverable"]
CMD ["serve"]
