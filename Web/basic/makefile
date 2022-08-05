version ?= 0.0.0

upload:
	ifndef APP_BUCKET
	$(error APP_BUCKET is not set)
	endif
	GOOS=linux GOARCH=amd64 go build -o web-basic-${version} .
	gsutil cp web-basic-${version} gs://${APP_BUCKET}/web-basic/web-basic-${version}
	rm web-basic-${version}

docker-start:
	docker build -t basic .
	docker run --name basic --env TARGET=aaa -p 8080:8080 basic

docker-stop:
	docker stop basic
	docker rm basic
