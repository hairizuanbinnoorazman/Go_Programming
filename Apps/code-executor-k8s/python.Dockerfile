FROM python:3.9-alpine
RUN adduser -D -u 3000 -g 3000 executor
USER executor