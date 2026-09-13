FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/proxy-cors-misconfiguration/ ./challenges/proxy-cors-misconfiguration/

WORKDIR /app/challenges/proxy-cors-misconfiguration
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /proxy-cors-misconfiguration .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /proxy-cors-misconfiguration /usr/bin/proxy-cors-misconfiguration

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["proxy-cors-misconfiguration"]
CMD ["serve"]
