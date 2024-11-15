FROM golang:1.22-alpine3.19
RUN adduser -D -u 3000 -g 3000 executor
USER executor