FROM ruby:3.3.6-alpine3.20
RUN adduser -D -u 3000 -g 3000 executor
USER executor