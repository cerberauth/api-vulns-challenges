FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o /jwt-weak-hmac-secret .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-weak-hmac-secret /usr/bin/jwt-weak-hmac-secret

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-weak-hmac-secret", "serve"]
CMD ["jwt-weak-hmac-secret"]
