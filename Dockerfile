FROM golang:1.14.3-alpine3.11 as builder

ENV GO111MODULE=on

WORKDIR /go/src/github.com/samuel-sujith/mwatch
COPY . .

RUN CGO_ENABLED=0 go build  -o /adapter github.com/samuel-sujith/mwatch

FROM alpine:3.11
RUN apk update \
    && apk add ca-certificates \
    && rm -rf /var/cache/apk/* \
    && update-ca-certificates

ENTRYPOINT ["/adapter", "--logtostderr=true"]
COPY --from=builder /adapter /