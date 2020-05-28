FROM golang:1.14
ADD . .
RUN go build -o app .
EXPOSE 8080
CMD ["./app"]