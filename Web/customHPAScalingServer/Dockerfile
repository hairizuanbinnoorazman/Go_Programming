FROM golang:1.15
WORKDIR /home
ADD go.mod go.sum ./
RUN go mod download
ADD main.go ./
RUN go build -o custom .