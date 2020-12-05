FROM golang:latest as builder
RUN mkdir -p /go/src/sk-kill
WORKDIR /go/src/sk-kill
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o kill .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/sk-kill/kill .
CMD ["./kill"]