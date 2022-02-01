start:
	APP_SECRET=aa API_KEY=aa go run main.go

docker-start:
	docker run -e APP_SECRET=aa -e API_KEY=aa -p 8080:8080 lol

check:
	curl -H "Authorization: aa" localhost:8080/api/lta-datamall/v1/bus-arrival?bus-stop-id=11111