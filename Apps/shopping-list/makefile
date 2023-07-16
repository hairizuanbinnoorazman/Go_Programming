project = $(shell gcloud config get project)
serviceurl = $(shell gcloud run services describe shopping-list --format=yaml --region=asia-southeast1 | yq '.status.address.url')
version = v0
item = new
id = xx

docker-prod:
	docker build -f prod.Dockerfile -t gcr.io/$(project)/shopping-list:$(version) .
	docker push gcr.io/$(project)/shopping-list:$(version)

docker-build:
	docker build -t gcr.io/$(project)/shopping-list:$(version) .

docker-push:
	docker push gcr.io/$(project)/shopping-list:$(version)

create:
	curl -X POST $(serviceurl)/api/shopping-list/v1/item -d '{"name": "$(item)"}'

list:
	curl -X GET $(serviceurl)/api/shopping-list/v1/item

delete:	
	curl -X DELETE $(serviceurl)/api/shopping-list/v1/item/$(id)

modify:
	curl $(serviceurl)/api/shopping-list/v1/item/$(id) -X PATCH -d '{"status":"in-cart"}'

format:
	cd ./frontend && \
	elm-format --yes src/

gen:
	cd ./frontend && \
	elm make --debug --output=main.js ./src/App.elm && \
	mv main.js ./files/main.js
	
gen-prod:
	cd ./frontend && \
	elm make --optimize --output=main.js ./src/App.elm && \
	uglifyjs main.js --compress "pure_funcs=[F2,F3,F4,F5,F6,F7,F8,F9,A2,A3,A4,A5,A6,A7,A8,A9],pure_getters,keep_fargs=false,unsafe_comps,unsafe" | uglifyjs --mangle --output main.min.js && \
	mv main.min.js ../static/main.js