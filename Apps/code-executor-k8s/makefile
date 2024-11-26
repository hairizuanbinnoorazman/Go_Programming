version?=v1
assessver?=v1
registry?=asia-southeast1-docker.pkg.dev/healthy-rarity-238313/code-assess/

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

upload-remote:
	go build -o app .
	docker build -t $(registry)lol:$(version) .
	docker push $(registry)lol:$(version)

upload-remote-executors:
	docker build -t $(registry)new-python:$(assessver) -f python.Dockerfile .
	docker build -t $(registry)new-golang:$(assessver) -f golang.Dockerfile .
	docker build -t $(registry)new-node:$(assessver) -f javascript.Dockerfile .
	docker build -t $(registry)new-ruby:$(assessver) -f ruby.Dockerfile .
	docker push $(registry)new-golang:$(assessver)
	docker push $(registry)new-python:$(assessver)
	docker push $(registry)new-node:$(assessver)
	docker push $(registry)new-ruby:$(assessver)