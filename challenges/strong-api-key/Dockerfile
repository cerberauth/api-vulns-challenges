FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o /strong-api-key .

FROM gcr.io/distroless/static-debian11:nonroot AS runner

WORKDIR /

COPY --from=builder --chown=nonroot:nonroot /strong-api-key /usr/bin/strong-api-key

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["strong-api-key"]
CMD ["serve"]
