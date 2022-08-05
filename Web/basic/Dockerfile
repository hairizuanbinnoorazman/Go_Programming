FROM golang:1.18
WORKDIR /helloworld
ADD . .
RUN go build -o app .
CMD ["/helloworld/app"]
EXPOSE 8080