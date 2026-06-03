FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-strong-eddsa-key/ ./challenges/jwt-strong-eddsa-key/

WORKDIR /app/challenges/jwt-strong-eddsa-key
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-strong-eddsa-key .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /app/challenges/jwt-strong-eddsa-key/keys /keys
COPY --from=builder --chown=nonroot:nonroot /jwt-strong-eddsa-key /usr/bin/jwt-strong-eddsa-key

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-strong-eddsa-key"]
CMD ["serve"]
