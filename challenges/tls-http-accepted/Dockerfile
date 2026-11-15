FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/tls-http-accepted/ ./challenges/tls-http-accepted/

WORKDIR /app/challenges/tls-http-accepted
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /tls-http-accepted .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /tls-http-accepted /usr/bin/tls-http-accepted

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["tls-http-accepted"]
CMD ["serve"]
