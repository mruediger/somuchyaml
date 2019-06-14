############################
FROM golang:alpine AS builder

# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/mruediger/somuchyaml/
COPY . .

RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/somuchyaml

############################
FROM scratch

COPY --from=builder /go/bin/somuchyaml /go/bin/somuchyaml
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY yaml.jpg .

ENTRYPOINT ["/go/bin/somuchyaml"]