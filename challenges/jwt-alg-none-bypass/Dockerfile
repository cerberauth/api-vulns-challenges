FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o /jwt-alg-none-bypass .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-alg-none-bypass /usr/bin/jwt-alg-none-bypass

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-alg-none-bypass"]
CMD ["serve"]
