FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/proxy-waf-bypass/ ./challenges/proxy-waf-bypass/

WORKDIR /app/challenges/proxy-waf-bypass
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /proxy-waf-bypass .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /proxy-waf-bypass /usr/bin/proxy-waf-bypass

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["proxy-waf-bypass"]
CMD ["serve"]
