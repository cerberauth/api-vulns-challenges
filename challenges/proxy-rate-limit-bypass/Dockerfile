FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/proxy-rate-limit-bypass/ ./challenges/proxy-rate-limit-bypass/

WORKDIR /app/challenges/proxy-rate-limit-bypass
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /proxy-rate-limit-bypass .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /proxy-rate-limit-bypass /usr/bin/proxy-rate-limit-bypass

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["proxy-rate-limit-bypass"]
CMD ["serve"]
