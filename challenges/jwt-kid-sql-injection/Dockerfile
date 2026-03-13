FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-kid-sql-injection/ ./challenges/jwt-kid-sql-injection/

RUN cd challenges/jwt-kid-sql-injection && CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-kid-sql-injection .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-kid-sql-injection /usr/bin/jwt-kid-sql-injection

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-kid-sql-injection"]
CMD ["serve"]
