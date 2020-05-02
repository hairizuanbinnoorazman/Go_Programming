FROM golang:1.14
WORKDIR /hahax
ADD go.mod go.sum ./
RUN go mod download
ADD main.go .
RUN go build -o app .
CMD ["./app"]