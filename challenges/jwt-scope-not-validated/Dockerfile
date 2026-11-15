FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-scope-not-validated/ ./challenges/jwt-scope-not-validated/

WORKDIR /app/challenges/jwt-scope-not-validated
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-scope-not-validated .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-scope-not-validated /usr/bin/jwt-scope-not-validated

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-scope-not-validated"]
CMD ["serve"]
