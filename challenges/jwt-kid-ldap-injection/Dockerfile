FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-kid-ldap-injection/ ./challenges/jwt-kid-ldap-injection/

WORKDIR /app/challenges/jwt-kid-ldap-injection
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-kid-ldap-injection .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-kid-ldap-injection /usr/bin/jwt-kid-ldap-injection

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-kid-ldap-injection"]
CMD ["serve"]
