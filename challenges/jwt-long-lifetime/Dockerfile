FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-long-lifetime/ ./challenges/jwt-long-lifetime/

WORKDIR /app/challenges/jwt-long-lifetime
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-long-lifetime .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-long-lifetime /usr/bin/jwt-long-lifetime

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-long-lifetime"]
CMD ["serve"]
