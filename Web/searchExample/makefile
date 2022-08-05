docker-start:
	docker build -t search .
	docker run -p 8080:8080 --name search search

docker-stop:
	docker stop search
	docker rm search