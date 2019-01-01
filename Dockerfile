FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
COPY . $GOPATH/src/github.com/pminnebach/velo
WORKDIR $GOPATH/src/github.com/pminnebach/velo
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/velo

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/velo /go/bin/velo
EXPOSE 8080
CMD ["/go/bin/velo"]
