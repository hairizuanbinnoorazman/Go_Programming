FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -trimpath -ldflags="-s -w" -o /out/loadlab-server ./cmd/server \
 && CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -trimpath -ldflags="-s -w" -o /out/loadlab-loadgen ./cmd/loadgen

FROM alpine:3.22
RUN apk add --no-cache ca-certificates && adduser -D -u 10001 loadlab
WORKDIR /app
COPY --from=build /out/loadlab-server /usr/local/bin/loadlab-server
COPY migrations ./migrations
USER loadlab
EXPOSE 8080
ENTRYPOINT ["loadlab-server"]
