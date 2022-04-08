FROM golang:1.14 as build
WORKDIR /app
ADD . .
RUN CGO_ENABLED=0 go build -o app .

FROM gcr.io/distroless/base-debian11:nonroot
COPY --from=build /app/app /app
EXPOSE 8080
CMD ["/app"] 