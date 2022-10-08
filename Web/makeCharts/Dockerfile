FROM golang:1.19 as source
WORKDIR /app
ADD go.mod go.sum .
RUN go mod download
ADD . .
RUN go build -o app .

FROM debian:buster
RUN apt update -y
RUN apt install -y chromium
COPY --from=source /app/app /app
CMD ["/app"]


