FROM golang:latest as builder
RUN mkdir -p /go/src/gateway-secondkill
WORKDIR /go/src/gateway-secondkill
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/gateway-secondkill/gateway .
CMD ["./gateway"]