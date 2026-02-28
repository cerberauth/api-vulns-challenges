FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-not-verified/ ./challenges/jwt-not-verified/

RUN cd challenges/jwt-not-verified && CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-not-verified .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /app/challenges/jwt-not-verified/keys /keys
COPY --from=builder --chown=nonroot:nonroot /jwt-not-verified /usr/bin/jwt-not-verified

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-not-verified"]
CMD ["serve"]
