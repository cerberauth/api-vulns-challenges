FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/cache-deception-normalization/ ./challenges/cache-deception-normalization/

WORKDIR /app/challenges/cache-deception-normalization
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /cache-deception-normalization .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /cache-deception-normalization /usr/bin/cache-deception-normalization

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["cache-deception-normalization"]
CMD ["serve"]
