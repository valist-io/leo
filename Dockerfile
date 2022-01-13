FROM golang:1.16
WORKDIR /go/src/github.com/valist-io/leo/
COPY . ./
ENV CGO_ENABLED 0
RUN go build ./cmd/leo

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/valist-io/leo/leo ./
ENTRYPOINT ["./leo"]
