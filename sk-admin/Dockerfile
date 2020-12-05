FROM golang:latest as builder
RUN mkdir -p /go/src/sk-admin
WORKDIR /go/src/sk-admin
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o admin .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/sk-admin/admin .
CMD ["./blog"]