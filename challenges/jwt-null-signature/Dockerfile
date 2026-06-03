FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-null-signature/ ./challenges/jwt-null-signature/

WORKDIR /app/challenges/jwt-null-signature
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-null-signature .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /app/challenges/jwt-null-signature/keys /keys
COPY --from=builder --chown=nonroot:nonroot /jwt-null-signature /usr/bin/jwt-null-signature

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-null-signature"]
CMD ["serve"]
