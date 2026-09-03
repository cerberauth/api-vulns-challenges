FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-claim-special-chars/ ./challenges/jwt-claim-special-chars/

WORKDIR /app/challenges/jwt-claim-special-chars
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-claim-special-chars .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-claim-special-chars /usr/bin/jwt-claim-special-chars

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-claim-special-chars"]
CMD ["serve"]
