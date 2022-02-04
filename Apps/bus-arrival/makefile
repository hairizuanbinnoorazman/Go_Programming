start:
	API_KEY=aa go run main.go

docker-start:
	docker run -e API_KEY=aa -p 8880:8080 lol

check:
	curl -H localhost:8080/api/lta-datamall/v1/bus-arrival?bus-stop-id=99999
