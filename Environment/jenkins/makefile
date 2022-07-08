build:
	docker build -t cjenkins .
run:
	docker run --name jenkins -p 8090:8080 -d cjenkins
stop:
	docker stop jenkins
	docker rm jenkins