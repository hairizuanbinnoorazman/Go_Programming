# Basic Example of Golang App talking via .sock

Note: This would only work on linux based systems (Both Host and Container is on Linux)

This step would create a container with a small app that creates a socket for file to connect to

```bash
docker build -t lol .
docker run -v $(pwd)/tmp:/tmp lol
```

To test that the app can accept data to the socket, you can run the following command. It would send data over to the socket and would print it out in the logs.

```bash
echo -e '\x66\x6f\x6f' | nc -U $(pwd)/tmp/go.sock
```

If you are testing this on linux OS-es, do make sure that you install `netcat-openbsd` in order to use the above `nc` command.

```
apt update && apt install -y netcat-openbsd
```

# Edge case

https://forums.docker.com/t/cant-connect-to-host-listening-unix-socket-from-container-vm/15526/2

If you try to build the main.go in a docker container and try to run it, you would face the issue where you would be unable to connect to it. Please refer to the link for additional details on this.

Apparently, docker.sock is a special socket file for this; sockets don't work out of the box between different OS-es