FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/proxy-info-disclosure/ ./challenges/proxy-info-disclosure/

WORKDIR /app/challenges/proxy-info-disclosure
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /proxy-info-disclosure .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /proxy-info-disclosure /usr/bin/proxy-info-disclosure

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["proxy-info-disclosure"]
CMD ["serve"]
