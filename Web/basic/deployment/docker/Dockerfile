FROM golang:1.18
WORKDIR /helloworld
ADD . .
RUN CGO_ENABLED=0 go build -o app .
CMD ["/helloworld/app"]
EXPOSE 8080