FROM golang:latest as builder
WORKDIR /go/src/github.com/pminnebach/velo
RUN go get -u -v -d github.com/go-resty/resty
RUN go get -u -v -d github.com/prometheus/client_golang/prometheus
RUN go get -u -v -d github.com/prometheus/client_golang/prometheus/promhttp
COPY main.go    .
RUN GOOS=linux GOARCH=arm go build -a -o velo .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/pminnebach/velo/velo .
EXPOSE 8080
CMD ["./velo"]
