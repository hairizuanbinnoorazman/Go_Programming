FROM golang:1.14
WORKDIR /home
ADD main.go .
RUN go build -o lol .
CMD ["/home/lol"]
EXPOSE 8080