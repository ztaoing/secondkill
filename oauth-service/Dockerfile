FROM golang:latest as builder
RUN mkdir -p /go/src/oauth
WORKDIR /go/src/oauth
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o oauth .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/oauth/oauth .
CMD ["./oauth"]