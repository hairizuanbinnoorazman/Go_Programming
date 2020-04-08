FROM golang
ADD . /go/src/github.com/hairizuanbinnoorazman/basicWithTracing
WORKDIR /go/src/github.com/hairizuanbinnoorazman/basicWithTracing
RUN go get
RUN go install github.com/hairizuanbinnoorazman/basicWithTracing
ENTRYPOINT ["/go/bin/basicWithTracing"]
EXPOSE 8080