grpc:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ticketing/ticketing.proto
build:
	go build -o lol ./server
docker-build:
	docker build --target client -t gcr.io/${GCP_PROJECT}/grpc-client:${VERSION} .
	docker build --target server -t gcr.io/${GCP_PROJECT}/grpc-server:${VERSION} .
docker-push:
	docker push gcr.io/${GCP_PROJECT}/grpc-server:${VERSION}
	docker push gcr.io/${GCP_PROJECT}/grpc-client:${VERSION}
deploy:
	kustomize edit set image grpc-client=gcr.io/${GCP_PROJECT}/grpc-client:${VERSION}
	kustomize edit set image grpc-server=gcr.io/${GCP_PROJECT}/grpc-server:${VERSION}
	kustomize build | kubectl apply -f -