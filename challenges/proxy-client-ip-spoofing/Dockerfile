FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/proxy-client-ip-spoofing/ ./challenges/proxy-client-ip-spoofing/

WORKDIR /app/challenges/proxy-client-ip-spoofing
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /proxy-client-ip-spoofing .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /proxy-client-ip-spoofing /usr/bin/proxy-client-ip-spoofing

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["proxy-client-ip-spoofing"]
CMD ["serve"]
