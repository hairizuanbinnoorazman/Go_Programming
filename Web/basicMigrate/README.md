# Basic Migrate

A basic golang app that embeds golang application that does migration as well as does simple store and gets from the database

## Migrations

First, we need to install the golang-migrate tool:

```bash
brew install golang-migrate
```

To create a migration sql script:

```bash
migrate create -dir migrations -seq -ext sql init_schema
```

We can test it on a mysql database

```bash
docker run --name some-mysql -e MYSQL_DATABASE=application -e MYSQL_ROOT_PASSWORD=my-secret-pw -e MYSQL_USER=user -e MYSQL_PASSWORD=password -p 3306:3306 -d mysql:5.7
```

Build the binary

```bash
go build -o lol .
```

Then, run migrate command and run the server

```bash
./lol migrate
./lol server
```

Run the curl command

```bash
curl -X POST localhost:8888/user -d '{"first_name": "test", "last_name":"test"}'
curl -X GET localhost:8888/user/1
```