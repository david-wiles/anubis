FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git
RUN apk --no-cache add ca-certificates

WORKDIR $GOPATH/src/anubis

COPY ./go.mod ./go.mod

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/anubis ./cmd/main.go

RUN chmod +x /go/bin/anubis

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/anubis /go/bin/anubis

ENTRYPOINT ["/go/bin/anubis"]
