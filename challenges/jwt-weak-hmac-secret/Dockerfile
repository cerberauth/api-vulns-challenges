FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-weak-hmac-secret/ ./challenges/jwt-weak-hmac-secret/

RUN cd challenges/jwt-weak-hmac-secret && CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-weak-hmac-secret .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-weak-hmac-secret /usr/bin/jwt-weak-hmac-secret

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-weak-hmac-secret"]
CMD ["serve"]
