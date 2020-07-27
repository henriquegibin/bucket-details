FROM golang:1.13.1-alpine3.10 AS builder

WORKDIR $GOPATH/src/github.com/henriquegibin/bucket-details

COPY . .

RUN apk update && apk add --no-cache git && go mod download \
  && CGO_ENABLED=0 GOOS=linux go build -o /plugin/bin/bucket-details \
  && chmod +x /plugin/bin/bucket-details

FROM scratch AS runner

WORKDIR /app

COPY --from=builder /plugin/bin/bucket-details .
COPY --from=builder /etc/ssl /etc/ssl

CMD ["/app/bucket-details"]
