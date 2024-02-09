# Builder stage
FROM golang:alpine AS builder

RUN apk update && apk add --no-cache \
  git \
  ca-certificates \
  tzdata \
  && update-ca-certificates

COPY .  $GOPATH/src/github.com/jcleira/encinitas-collector-go/
WORKDIR $GOPATH/src/github.com/jcleira/encinitas-collector-go/

RUN go build -mod=vendor -o /go/bin/encinitas-collector-go
RUN CGO_ENABLED=0 GOOS=linux \
  go build -a -installsuffix cgo -mod=vendor \
  -o /go/bin/encinitas-collector-go .

# Runner stage
FROM scratch

COPY --from=builder /go/bin/encinitas-collector-go /go/bin/encinitas-collector-g-
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ENTRYPOINT ["/go/bin/encinitas-collector-go"]

