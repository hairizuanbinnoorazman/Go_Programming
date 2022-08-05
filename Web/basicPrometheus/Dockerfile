FROM golang:1.15
WORKDIR /app
ADD go.mod go.sum .
RUN go mod download
ADD main.go .
RUN go build -o app .
CMD /app/app
EXPOSE 8080