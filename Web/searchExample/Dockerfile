FROM golang:1.18
WORKDIR /search-example
ADD . .
RUN go build -o app ./cmd
CMD ["/search-example/app"]
EXPOSE 8080