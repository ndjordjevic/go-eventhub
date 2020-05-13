up:
	docker-compose -f cmd/server-echo/docker-compose.yml up

down:
	docker-compose -f cmd/server-echo/docker-compose.yml down

build:
	docker build -t ndjordjevic/server-echo -f cmd/server-echo/Dockerfile .