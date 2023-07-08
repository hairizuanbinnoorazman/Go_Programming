FROM golang:1.19 as builder
WORKDIR /shopping-list
COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN go build -o app .

FROM debian:bookworm-slim
RUN apt update && apt install -y ca-certificates
WORKDIR /opt/shopping-list
# RUN useradd -ms /bin/bash shopping-list
COPY --from=builder /shopping-list/app /opt/shopping-list/app
# RUN chown shopping-list /opt/shopping-list/app
RUN chmod +x /opt/shopping-list/app
# USER shopping-list
CMD ["/opt/shopping-list/app"]