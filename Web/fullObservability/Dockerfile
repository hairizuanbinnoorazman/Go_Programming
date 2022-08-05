FROM golang:1.18
WORKDIR /full-observability
ADD go.mod go.sum .
RUN go mod download
ADD . .
RUN go build -o app .
CMD ["/full-observability/app"]
EXPOSE 8080