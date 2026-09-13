FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/proxy-template-injection/ ./challenges/proxy-template-injection/

WORKDIR /app/challenges/proxy-template-injection
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /proxy-template-injection .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /proxy-template-injection /usr/bin/proxy-template-injection

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["proxy-template-injection"]
CMD ["serve"]
