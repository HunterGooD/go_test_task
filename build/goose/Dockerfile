FROM golang:1.22.6-alpine3.20 AS builder
RUN apk update && apk --no-cache add build-base git
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:3.20
COPY --from=builder /go/bin/goose /go/bin/goose

WORKDIR /app
# COPY ./db/migrations /app/db/migrations
# ENTRYPOINT ["/go/bin/goose", "-dir=/app/db/migrations/", "-v"]
# CMD [ "up" ]
CMD ["/go/bin/goose", "-dir=/app/db/migrations/", "-v", "up"]
