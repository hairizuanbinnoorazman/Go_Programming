FROM golang:1.21 as builder
WORKDIR /headlessservice
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o app .

FROM debian:bullseye-slim as prod
RUN apt update && apt install -y ca-certificates
COPY --from=builder /headlessservice/app /usr/bin/app
WORKDIR /usr/bin
EXPOSE 8080
CMD ["app"]