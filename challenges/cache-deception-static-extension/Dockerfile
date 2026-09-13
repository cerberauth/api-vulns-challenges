FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/cache-deception-static-extension/ ./challenges/cache-deception-static-extension/

WORKDIR /app/challenges/cache-deception-static-extension
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /cache-deception-static-extension .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /cache-deception-static-extension /usr/bin/cache-deception-static-extension

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["cache-deception-static-extension"]
CMD ["serve"]
