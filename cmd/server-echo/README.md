go run -race main.go

go run cmd/client-gorilla/main.go -user ndjordjevic -addr localhost:8080

go run cmd/client-gorilla/main.go -user vpopovic -addr localhost:8080

go run cmd/nats-client/main.go -s localhost sb-events "ndjordjevic,new_order123"   

docker pull nats:latest   
docker run --name nats -p 4222:4222 --network event-hub-net -ti nats:latest

docker build -t ndjordjevic/server-echo -f cmd/server-echo/Dockerfile . (from the repo root)

docker run --name event-hub -tid --network event-hub-net -p 8080:8080 -e NATS_ADDR=localhost ndjordjevic/server-echo
