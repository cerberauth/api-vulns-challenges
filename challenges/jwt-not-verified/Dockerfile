FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o /jwt-not-verified .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-not-verified /usr/bin/jwt-not-verified

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-not-verified"]
CMD ["serve"]
