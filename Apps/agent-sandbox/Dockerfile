FROM golang:1.26.1-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/agent-sandbox-demo .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/agent-sandbox-demo /agent-sandbox-demo
EXPOSE 8080
ENTRYPOINT ["/agent-sandbox-demo"]
