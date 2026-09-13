FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/cache-conditional-validators/ ./challenges/cache-conditional-validators/

WORKDIR /app/challenges/cache-conditional-validators
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /cache-conditional-validators .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /cache-conditional-validators /usr/bin/cache-conditional-validators

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["cache-conditional-validators"]
CMD ["serve"]
