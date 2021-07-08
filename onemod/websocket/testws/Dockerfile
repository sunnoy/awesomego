FROM golang:1.16.3

WORKDIR /root/
COPY server.go .
RUN go get github.com/gorilla/websocket
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /root/server .
CMD ["./server"]