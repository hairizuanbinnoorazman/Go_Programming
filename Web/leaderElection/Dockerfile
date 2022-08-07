FROM golang:1.18
WORKDIR /leader-election
ADD go.mod go.sum .
RUN go mod download
ADD . .
RUN go build -o app .
CMD ["/leader-election/app"]
EXPOSE 8080