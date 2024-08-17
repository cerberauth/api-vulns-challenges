FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o /jwt-null-signature .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /jwt-null-signature /usr/bin/jwt-null-signature

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-null-signature", "serve"]
CMD ["jwt-null-signature"]
