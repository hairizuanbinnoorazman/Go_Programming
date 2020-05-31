FROM golang
WORKDIR /go/src/github.com/hairizuanbinnoorazman/basicwithstackdriver
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN go install -ldflags="-X 'main.ServiceName=application' -X 'main.Version=0.1.5'" github.com/hairizuanbinnoorazman/basicwithstackdriver
CMD ["/go/bin/basicwithstackdriver"]
EXPOSE 8080