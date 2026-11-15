FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/auth-unauthenticated-endpoint/ ./challenges/auth-unauthenticated-endpoint/

WORKDIR /app/challenges/auth-unauthenticated-endpoint
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /auth-unauthenticated-endpoint .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /auth-unauthenticated-endpoint /usr/bin/auth-unauthenticated-endpoint

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["auth-unauthenticated-endpoint"]
CMD ["serve"]
