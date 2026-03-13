FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-kid-path-traversal/ ./challenges/jwt-kid-path-traversal/

RUN cd challenges/jwt-kid-path-traversal && CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-kid-path-traversal .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /app/challenges/jwt-kid-path-traversal/keys /keys
COPY --from=builder --chown=nonroot:nonroot /jwt-kid-path-traversal /usr/bin/jwt-kid-path-traversal

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-kid-path-traversal"]
CMD ["serve"]
