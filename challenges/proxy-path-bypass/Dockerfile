FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/proxy-path-bypass/ ./challenges/proxy-path-bypass/

WORKDIR /app/challenges/proxy-path-bypass
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /proxy-path-bypass .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /proxy-path-bypass /usr/bin/proxy-path-bypass

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["proxy-path-bypass"]
CMD ["serve"]
