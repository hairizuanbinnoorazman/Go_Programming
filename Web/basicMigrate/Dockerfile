FROM golang:1.20 as builder
WORKDIR /basicmigrate
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o app .

FROM debian:bookworm-slim
RUN apt update && \
    apt install -y ca-certificates && \
    apt clean && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /basicmigrate
COPY --from=builder /basicmigrate/app /usr/local/bin/app
CMD ["app"]
EXPOSE 8080
