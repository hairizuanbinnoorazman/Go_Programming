FROM golang:1.18 as builder
WORKDIR /helloworld
ADD . .
RUN CGO_ENABLED=0 go build -o app .

FROM alpine:3.16
RUN apk add gcompat
WORKDIR /helloworld
COPY --from=builder /helloworld/app /helloworld/app
CMD ["/helloworld/app"]
EXPOSE 8080

