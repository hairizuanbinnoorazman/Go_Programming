peer1:
	BT_PORT=6181 HTTP_PORT=8080 PEER_URL=localhost:6182 go run main.go

peer2:
	FOLDER_PATH=./zzz BT_PORT=6183 HTTP_PORT=8082 PEER_URL=localhost:6181 go run main.go

addtorrent:
	curl http://localhost:8080/add-torrent

addtorrent2:
	curl http://localhost:8082/add-torrent

starttorrent2:
	curl http://localhost:8082/start

allstart: addtorrent addtorrent2 starttorrent2
