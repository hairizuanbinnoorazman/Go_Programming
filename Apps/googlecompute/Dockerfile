FROM golang:1.23 as builder
WORKDIR /googlecompute
ADD go.mod go.sum ./
RUN go mod download
ADD main.go .
RUN go build -o app .

FROM debian:bookworm-slim
RUN apt update && apt install -y ca-certificates
WORKDIR /opt/googlecompute
COPY --from=builder /googlecompute/app /opt/googlecompute/app
CMD ["/opt/googlecompute/app"]