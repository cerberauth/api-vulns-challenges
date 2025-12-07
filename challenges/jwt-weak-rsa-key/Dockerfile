FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o /jwt-weak-rsa-key .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /app/keys /keys
COPY --from=builder --chown=nonroot:nonroot /jwt-weak-rsa-key /usr/bin/jwt-weak-rsa-key

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["jwt-weak-rsa-key"]
CMD ["serve"]
