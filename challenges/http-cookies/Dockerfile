FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o /http-cookies .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /http-cookies /usr/bin/http-cookies

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["http-cookies", "serve"]
CMD ["http-cookies"]
