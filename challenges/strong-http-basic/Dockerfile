FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/strong-http-basic/ ./challenges/strong-http-basic/

RUN cd challenges/strong-http-basic && CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /strong-http-basic .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /strong-http-basic /usr/bin/strong-http-basic

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["strong-http-basic"]
CMD ["serve"]
