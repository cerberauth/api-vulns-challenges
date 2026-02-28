FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/http-misconfigurations/ ./challenges/http-misconfigurations/

RUN cd challenges/http-misconfigurations && CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /http-misconfigurations .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /http-misconfigurations /usr/bin/http-misconfigurations

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["http-misconfigurations"]
CMD ["serve"]
