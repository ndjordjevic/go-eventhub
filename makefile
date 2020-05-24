up:
	docker-compose -f cmd/server-echo/docker-compose.yml up

down:
	docker-compose -f cmd/server-echo/docker-compose.yml down

build:
	docker build -t ndjordjevic/server-echo -f cmd/server-echo/Dockerfile .

proto:
	protoc --proto_path=api --go_out=./internal/protogen/api --go_opt=paths=source_relative ./api/instrument.proto
