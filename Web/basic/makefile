version ?= 0.0.0

ifndef APP_BUCKET
$(error APP_BUCKET is not set)
endif

upload:
	GOOS=linux GOARCH=amd64 go build -o web-basic-${version} .
	gsutil cp web-basic-${version} gs://${APP_BUCKET}/web-basic/web-basic-${version}
	rm web-basic-${version}