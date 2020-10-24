# Basic application

A simple Golang application that interacts with a mysql database to create a records

# Local development

```
docker run \
    -p 3306:3306 \
    -e MYSQL_ROOT_PASSWORD=root \
    -e MYSQL_DATABASE=testmysql \
    -e MYSQL_USER=username \
    -e MYSQL_PASSWORD=password \
    mysql:5.7
```