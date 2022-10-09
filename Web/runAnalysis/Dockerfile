FROM golang:1.19 as source
WORKDIR /app
ADD go.mod go.sum .
RUN go mod download
ADD . .
RUN go build -o app .

FROM debian:buster-slim
RUN apt-get update \
    && apt-get install -y --force-yes ca-certificates \
    && apt-get clean \
    && apt-get autoremove \
    && rm -rf /var/lib/apt/lists/*
COPY --from=source /app/app /app
CMD ["/app"]
