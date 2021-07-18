FROM golang:1.16 as base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM base as client
COPY . .
RUN go build -o app ./client
CMD ["/app/app"]

FROM base as server
COPY . .
RUN go build -o app ./server
CMD ["/app/app"]
EXPOSE 12345