FROM golang:1.19 as builder
WORKDIR /basic
COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN go build -o app ./cmd/backend

FROM ubuntu:20.04 as prod
RUN apt update && apt install -y ca-certificates
COPY --from=builder /basic/app /usr/bin/app
WORKDIR /usr/bin
EXPOSE 8080
CMD ["app", "server"]