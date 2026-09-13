FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/proxy-http2-authority-spoofing/ ./challenges/proxy-http2-authority-spoofing/

WORKDIR /app/challenges/proxy-http2-authority-spoofing
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /proxy-http2-authority-spoofing .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /proxy-http2-authority-spoofing /usr/bin/proxy-http2-authority-spoofing

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["proxy-http2-authority-spoofing"]
CMD ["serve"]
