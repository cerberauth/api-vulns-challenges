FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/proxy-request-smuggling/ ./challenges/proxy-request-smuggling/

WORKDIR /app/challenges/proxy-request-smuggling
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /proxy-request-smuggling .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /proxy-request-smuggling /usr/bin/proxy-request-smuggling

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["proxy-request-smuggling"]
CMD ["serve"]
