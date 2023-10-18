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
sudo apt install -y maria-server maria-client
```

Then enter into mysql command line

```bash
mysql
```

Then create the database and create user

```
CREATE DATABASE testmysql
CREATE USER username IDENTIFIED BY 'password'
```

Copy the binary over

```bash
sftp recordmaker hairizuan@<ip address>:/root/recordmaker
```