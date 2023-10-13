FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o /jwt-strong-eddsa-key .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-strong-eddsa-key /usr/bin/jwt-strong-eddsa-key

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-strong-eddsa-key", "serve"]
CMD ["jwt-strong-eddsa-key"]
