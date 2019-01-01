FROM golang:1.11-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
WORKDIR /go/src/github.com/ricardocampos/goauth
COPY . .

RUN go get -d -v ./...
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/goauth

FROM scratch
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/src/github.com/ricardocampos/goauth /go/src/github.com/ricardocampos/goauth
COPY --from=builder /go/bin/goauth /go/bin/goauth

WORKDIR /go/src/github.com/ricardocampos/goauth
USER nobody
EXPOSE 8080
ENTRYPOINT ["/go/bin/goauth"]