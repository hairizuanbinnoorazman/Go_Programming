FROM golang:1.13
RUN apt update && apt install -y netcat-openbsd
ADD . .
RUN go build -o app ./main.go
CMD ["/go/app"]