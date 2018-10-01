FROM golang:alpine as builder
RUN apk add --no-cache git
COPY . $GOPATH/src/emon/
WORKDIR $GOPATH/src/emon/
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/emon

FROM scratch
COPY --from=builder /go/bin/emon /go/bin/emon
EXPOSE 8113
ENTRYPOINT ["/go/bin/emon"]
