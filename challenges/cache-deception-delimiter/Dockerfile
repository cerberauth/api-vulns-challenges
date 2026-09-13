FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/cache-deception-delimiter/ ./challenges/cache-deception-delimiter/

WORKDIR /app/challenges/cache-deception-delimiter
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /cache-deception-delimiter .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /cache-deception-delimiter /usr/bin/cache-deception-delimiter

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["cache-deception-delimiter"]
CMD ["serve"]
