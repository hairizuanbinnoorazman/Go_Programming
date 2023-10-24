FROM golang:1.18 as builder
WORKDIR /helloworld
ADD . .
RUN CGO_ENABLED=0 go build -o app .

FROM debian:bookworm-slim
RUN apt update && \
    apt install -y ca-certificates && \
    apt clean && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /helloworld
COPY --from=builder /helloworld/app /helloworld/app
CMD ["/helloworld/app"]
EXPOSE 8080
