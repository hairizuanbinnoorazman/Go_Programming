FROM golang
ADD . /go/src/github.com/knative/docs/helloworld
RUN go install github.com/knative/docs/helloworld
ENTRYPOINT /go/bin/helloworld
EXPOSE 8080