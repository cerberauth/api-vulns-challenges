FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o /strong-http-basic .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /strong-http-basic /usr/bin/strong-http-basic

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["strong-http-basic"]
CMD ["serve"]
