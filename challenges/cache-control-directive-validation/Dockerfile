FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/cache-control-directive-validation/ ./challenges/cache-control-directive-validation/

WORKDIR /app/challenges/cache-control-directive-validation
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /cache-control-directive-validation .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /cache-control-directive-validation /usr/bin/cache-control-directive-validation

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["cache-control-directive-validation"]
CMD ["serve"]
