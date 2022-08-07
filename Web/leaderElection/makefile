deploy:
	kustomize build ./deployment | kubectl apply -f -

deploy-debug:
	kubectl create deployment debug --image=nginx

delete:
	kubectl delete lease app-lock
	kubectl delete configmap app-lock
	kustomize build ./deployment | kubectl delete -f -

build:
	docker build -t gcr.io/XXX/leader-election:0.0.12 .
	docker push gcr.io/XXX/leader-election:0.0.12
	