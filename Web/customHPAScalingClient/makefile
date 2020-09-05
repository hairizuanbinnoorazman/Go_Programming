deploy:
	kustomize build . | kubectl apply -f -

drop:
	kubectl delete -k .