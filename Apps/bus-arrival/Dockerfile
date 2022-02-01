FROM golang:1.17 as builder
WORKDIR /bus-arrival
COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN go build -o app .

FROM debian:bookworm-slim
WORKDIR /opt/bus-arrival
RUN useradd -ms /bin/bash bus-arrival
COPY --from=builder /bus-arrival/app /opt/bus-arrival/app
RUN chown bus-arrival /opt/bus-arrival/app
RUN chmod +x /opt/bus-arrival/app
USER bus-arrival
CMD ["/opt/bus-arrival/app"]