FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/auth-not-verified/ ./challenges/auth-not-verified/

RUN cd challenges/auth-not-verified && CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /auth-not-verified .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /auth-not-verified /usr/bin/auth-not-verified

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["auth-not-verified"]
CMD ["serve"]
