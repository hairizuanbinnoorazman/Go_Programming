FROM golang
ADD . /go/src/github.com/hairizuanbinnoorazman/full-observability
WORKDIR /go/src/github.com/hairizuanbinnoorazman/full-observability
RUN go get
RUN go install github.com/hairizuanbinnoorazman/full-observability
ENTRYPOINT ["/go/bin/full-observability"]
EXPOSE 8080 6060