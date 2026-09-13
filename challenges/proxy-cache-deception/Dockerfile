FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/proxy-cache-deception/ ./challenges/proxy-cache-deception/

WORKDIR /app/challenges/proxy-cache-deception
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /proxy-cache-deception .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /proxy-cache-deception /usr/bin/proxy-cache-deception

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["proxy-cache-deception"]
CMD ["serve"]
