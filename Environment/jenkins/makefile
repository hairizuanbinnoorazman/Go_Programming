build:
	docker build -t cjenkins .
run:
	docker run --name jenkins -p 8090:8080 --network cicd -d cjenkins
stop:
	docker stop jenkins
	docker rm jenkins
run-agent:
	docker run --name agent -d --network cicd jenkins/agent java -jar /usr/share/jenkins/agent.jar -url http://jenkins:8080/ -secret xxx -name zzz -workDir "/home/jenkins"
stop-agent:
	docker stop agent
	docker rm agent