FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-alg-none-bypass/ ./challenges/jwt-alg-none-bypass/

RUN cd challenges/jwt-alg-none-bypass && CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-alg-none-bypass .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-alg-none-bypass /usr/bin/jwt-alg-none-bypass

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-alg-none-bypass"]
CMD ["serve"]
