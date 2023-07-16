project = $(shell gcloud config get project)
serviceurl = $(shell gcloud run services describe shopping-list --format=yaml --region=asia-southeast1 | yq '.status.address.url')
version = v0
git_hash=$$(git rev-parse --short HEAD)
image_repo?=""
image_version?=latest

build-prod-bin: gen-prod
	GOOS=linux GOARCH=amd64 go build -trimpath -tags=embedfrontend -o ./cmd/backend/app ./cmd/backend

build-images: 
	docker build -t $(image_repo)backend:$(image_version) ./cmd/backend

stack-up:
	cd deployment/docker-compose && docker-compose up

stack-down:
	cd deployment/docker-compose && docker-compose down

reup: stack-down stack-up

format:
	cd ./cmd/frontend && \
	elm-format --yes src/

gen:
	cd ./cmd/frontend && \
	elm make --debug --output=main.js ./src/App.elm && \
	mv main.js ./files/main.js

gen-prod:
	cd ./cmd/frontend && \
	elm make --optimize --output=main.js ./src/App.elm && \
	uglifyjs main.js --compress "pure_funcs=[F2,F3,F4,F5,F6,F7,F8,F9,A2,A3,A4,A5,A6,A7,A8,A9],pure_getters,keep_fargs=false,unsafe_comps,unsafe" | uglifyjs --mangle --output main.min.js && \
	mv main.min.js ../backend/static/main.js