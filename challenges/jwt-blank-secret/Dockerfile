FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-blank-secret/ ./challenges/jwt-blank-secret/

WORKDIR /app/challenges/jwt-blank-secret
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-blank-secret .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-blank-secret /usr/bin/jwt-blank-secret

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-blank-secret"]
CMD ["serve"]
