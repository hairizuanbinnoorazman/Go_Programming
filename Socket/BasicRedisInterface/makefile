redis:
	docker run --name=redis -d -p 6379:6379 redis
stop:
	docker stop redis
	docker rm redis
fakeserver:
	./lol
fakeserver-compat:
	ENABLE_HELLO_COMMAND=true ENABLE_BROKEN_UP_COMMANDS=true ./lol