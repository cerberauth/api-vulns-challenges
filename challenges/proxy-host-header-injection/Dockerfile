FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/proxy-host-header-injection/ ./challenges/proxy-host-header-injection/

WORKDIR /app/challenges/proxy-host-header-injection
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /proxy-host-header-injection .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /proxy-host-header-injection /usr/bin/proxy-host-header-injection

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["proxy-host-header-injection"]
CMD ["serve"]
