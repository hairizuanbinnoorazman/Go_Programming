# Using Golang with Redis

A very basic example of how to use redis with Golang

## Quick start

Run the following docker container and expose it accordingly:  
Ref: https://hub.docker.com/_/redis

```bash
docker run --name some-redis -p 6379:6379 -d redis
```


Setting foo value
set foo zzz ex 20: OK
Doing a ping command
Doing a get command
Value of foo: zzz
Value of foo: get foo: zzz