FROM golang:1.16
WORKDIR /app
ADD main.go .
RUN go build -o lol main.go
CMD ["/app/lol"]