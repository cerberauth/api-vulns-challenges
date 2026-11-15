FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-kid-command-injection/ ./challenges/jwt-kid-command-injection/

WORKDIR /app/challenges/jwt-kid-command-injection
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-kid-command-injection .

FROM alpine:3.22 AS runner

WORKDIR /

COPY --from=builder /jwt-kid-command-injection /usr/bin/jwt-kid-command-injection

EXPOSE 8080

USER nobody:nobody

ENTRYPOINT ["jwt-kid-command-injection"]
CMD ["serve"]
