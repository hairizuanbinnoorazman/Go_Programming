FROM golang:1.21 as builder
WORKDIR /helloworld
COPY . .
RUN CGO_ENABLED=0 go build -o app ./cmd/auth

FROM debian:bookworm-slim
RUN apt update && \
    apt install -y ca-certificates && \
    apt clean && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /helloworld
COPY --from=builder /helloworld/app /helloworld/app
COPY ./cmd/auth/layout.html /helloworld/layout.html
CMD ["/helloworld/app"]
EXPOSE 8080
