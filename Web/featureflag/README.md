## Feature flag

2 examples of applications provided of how we can set it up applications with feature flags. There are numerous way to implement this:

There are 2 ways that has been implemented here:

- Use environment variables as a way as feature flag
- Use feature platform tool (unleash) to handle feature flags for an application

## Quick start

Default username and password for unleash-server: `admin:unleash4all`

## Deploying on a GCE instance

- Create VM and expose both port 4242 and port 80
- Install docker on the GCE instance
- Install nginx on the GCE instance
- Run docker-compose.yaml on the GCE instance

```bash
# To allow current user to use docker commands without sudo
sudo usermod -aG docker $USER

# Create app user
sudo useradd app
sudo groupadd app

# Store the systemd file here: /lib/systemd/system/app.service

# Building the binary
GOOS=linux CGO_ENABLED=0 go build -o app ./cmd/unleash/
scp app <user>@<hostname>:app
sudo mv app /usr/local/bin/app
```