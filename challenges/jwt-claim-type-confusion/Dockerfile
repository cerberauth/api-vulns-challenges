FROM golang:1.26 AS builder

WORKDIR /app

COPY common/ ./common/
COPY challenges/jwt-claim-type-confusion/ ./challenges/jwt-claim-type-confusion/

WORKDIR /app/challenges/jwt-claim-type-confusion
RUN CGO_ENABLED=0 GOWORK=off GOOS=linux go build -o /jwt-claim-type-confusion .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-claim-type-confusion /usr/bin/jwt-claim-type-confusion

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-claim-type-confusion"]
CMD ["serve"]
