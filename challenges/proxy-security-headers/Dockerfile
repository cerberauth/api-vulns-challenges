FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/proxy-security-headers/ ./challenges/proxy-security-headers/

WORKDIR /app/challenges/proxy-security-headers
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /proxy-security-headers .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /proxy-security-headers /usr/bin/proxy-security-headers

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["proxy-security-headers"]
CMD ["serve"]
