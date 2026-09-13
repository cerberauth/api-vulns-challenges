FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/proxy-open-ssrf/ ./challenges/proxy-open-ssrf/

WORKDIR /app/challenges/proxy-open-ssrf
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /proxy-open-ssrf .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /proxy-open-ssrf /usr/bin/proxy-open-ssrf

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["proxy-open-ssrf"]
CMD ["serve"]
