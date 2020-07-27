FROM golang:1.13.1-alpine3.10 AS builder

WORKDIR $GOPATH/src/github.com/henriquegibin/bucket-details

COPY . .

RUN apk update && apk add --no-cache git build-base && go mod download \
  && CGO_ENABLED=0 GOOS=linux go build -o /app/bin/bucket-details \
  && chmod +x /app/bin/bucket-details

FROM builder as tests
RUN go test -cover ./src/aws/ \
  && go test -cover ./src/errorcheck \
  && go test -cover ./src/genericfunctions

FROM scratch AS runner

WORKDIR /app

COPY --from=builder /app/bin/bucket-details .
COPY --from=builder /etc/ssl /etc/ssl

CMD ["/app/bucket-details"]
