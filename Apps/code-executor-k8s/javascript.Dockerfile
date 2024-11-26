FROM node:alpine3.20
RUN adduser -D -u 3000 -g 3000 executor
USER executor