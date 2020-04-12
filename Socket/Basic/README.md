# Basic Example of Golang App talking via .sock

```bash
echo -e '\x66\x6f\x6f' | nc -U $(pwd)/tmp/go.sock
```

```bash
docker build -t lol .
docker run -v $(pwd)/tmp:/tmp lol
```

```
apt update && apt install -y netcat-openbsd
```