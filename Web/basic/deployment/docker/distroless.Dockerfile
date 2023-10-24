FROM golang:1.18 as builder
WORKDIR /helloworld
ADD . .
RUN CGO_ENABLED=0 go build -o app .

FROM gcr.io/distroless/static-debian11:nonroot
WORKDIR /helloworld
COPY --from=builder /helloworld/app /helloworld/app
CMD ["/helloworld/app"]
EXPOSE 8080

