FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/strong-api-key/ ./challenges/strong-api-key/

WORKDIR /app/challenges/strong-api-key
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /strong-api-key .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /strong-api-key /usr/bin/strong-api-key

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["strong-api-key"]
CMD ["serve"]
