# Using Golang with Memcached

A very basic example of how to use memcached with Golang

## Quick start

Run the following docker container and expose it accordingly:  
Ref: https://hub.docker.com/_/memcached

```bash
docker run --name my-memcache -p 11211:11211 -d memcached:1.6 memcached -m 64
```
