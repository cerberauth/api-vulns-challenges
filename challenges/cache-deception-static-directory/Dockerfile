FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/cache-deception-static-directory/ ./challenges/cache-deception-static-directory/

WORKDIR /app/challenges/cache-deception-static-directory
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /cache-deception-static-directory .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /cache-deception-static-directory /usr/bin/cache-deception-static-directory

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["cache-deception-static-directory"]
CMD ["serve"]
