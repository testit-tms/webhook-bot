FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /webhook-bot ./cmd/webhook-bot



FROM alpine:latest

ENV CONFIG_PATH=/app/config.yaml

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /webhook-bot .

COPY ./config/config.yaml .

EXPOSE 8080

CMD ["./webhook-bot"]
