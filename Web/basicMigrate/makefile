
all-mysql: start-mysql build
	sleep 30
	make migrate-mysql 
	make start-app-mysql
all-sqlite: build create-sqlite migrate-sqlite start-app-sqlite

start-mysql:
	docker run --name some-mysql -e MYSQL_DATABASE=application -e MYSQL_ROOT_PASSWORD=my-secret-pw -e MYSQL_USER=user -e MYSQL_PASSWORD=password -p 3306:3306 -d mysql:5.7

stop-mysql:
	docker stop some-mysql
	docker rm some-mysql

build:
	go build -o lol .

migrate-mysql:
	DATABASE_NAME=application DATABASE_USER=user DATABASE_PASSWORD=password \
	DATABASE_HOST=localhost DATABASE_TYPE=mysql \
	./lol migrate

start-app-mysql:
	DATABASE_NAME=application DATABASE_USER=user DATABASE_PASSWORD=password \
	DATABASE_HOST=localhost DATABASE_TYPE=mysql \
	./lol server

migrate-sqlite:
	DATABASE_NAME=application DATABASE_USER=user DATABASE_PASSWORD=password \
	DATABASE_HOST=localhost DATABASE_TYPE=sqlite SQLITE_FILE=application.db \
	./lol migrate

start-app-sqlite:
	DATABASE_NAME=application DATABASE_USER=user DATABASE_PASSWORD=password \
	DATABASE_HOST=localhost DATABASE_TYPE=sqlite SQLITE_FILE=application.db \
	./lol server

test-app:
	curl -X GET localhost:8888/health
	curl -X POST localhost:8888/user -d '{"first_name":"zzz","last_name":"zzz"}'
	curl -X GET localhost:8888/user/1

create-sqlite:
	sqlite3 application.db ".databases"

