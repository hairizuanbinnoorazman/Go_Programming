version?=v1

deploy:
	go build -o app .
	docker build -t lol:$(version) .
	minikube image load lol:$(version)
	 yq -i '.images[0].newTag = "$(version)"' kustomization.yaml
	kubectl apply -k .
	
upload:
	docker build -t new-python:$(version) -f python.Dockerfile .
	minikube image load new-python:$(version)
	# docker build -t new-golang:$(version) -f golang.Dockerfile .
	# minikube image load new-golang:$(version)