FROM golang:latest as builder
RUN mkdir -p /go/src/sk-app
WORKDIR /go/src/sk-app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/sk-app/app .
CMD ["./app"]