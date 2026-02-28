FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-weak-rsa-key/ ./challenges/jwt-weak-rsa-key/

RUN cd challenges/jwt-weak-rsa-key && CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-weak-rsa-key .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /app/challenges/jwt-weak-rsa-key/keys /keys
COPY --from=builder --chown=nonroot:nonroot /jwt-weak-rsa-key /usr/bin/jwt-weak-rsa-key

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-weak-rsa-key"]
CMD ["serve"]
