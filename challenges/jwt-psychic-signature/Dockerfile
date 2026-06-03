FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-psychic-signature/ ./challenges/jwt-psychic-signature/

WORKDIR /app/challenges/jwt-psychic-signature
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-psychic-signature .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /app/challenges/jwt-psychic-signature/keys /keys
COPY --from=builder --chown=nonroot:nonroot /jwt-psychic-signature /usr/bin/jwt-psychic-signature

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-psychic-signature"]
CMD ["serve"]
