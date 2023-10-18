# Basic application

A simple Golang application that interacts with a mysql database to create a records

## Local development

```
docker run \
    -p 3306:3306 \
    -e MYSQL_ROOT_PASSWORD=root \
    -e MYSQL_DATABASE=testmysql \
    -e MYSQL_USER=username \
    -e MYSQL_PASSWORD=password \
    mysql:5.7
```

## Fresh virtual machine

First, install mariadb - reason for using mariadb install of mysql is due to difficult to install mysql simply.

```bash
sudo apt update
sudo apt install -y mariadb-server mariadb-client
```

Then enter into mysql command line

```bash
mysql
```

Then create the database and create user

```
CREATE DATABASE testmysql;
CREATE USER username IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON `testmysql`.* TO 'username';
```

Build the binary

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o recordmaker .
```

Copy the binary over

```bash
scp recordmaker hairizuan@<ip address>:/home/hairizuan/recordmaker
```

Mysql command to check records on database:

```
select * from `testmysql`.`users` order by `updated_at` desc limit 10 ;
```