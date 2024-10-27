version?=v1

deploy:
	go build -o app .
	docker build -t lol:$(version) .
	minikube image load lol:$(version)
	 yq -i '.images[0].newTag = "$(version)"' kustomization.yaml
	kubectl apply -k .
	