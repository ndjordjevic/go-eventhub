go run -race main.go

go run cmd/client-gorilla/main.go -user ndjordjevic -addr localhost:8080

go run cmd/client-gorilla/main.go -user vpopovic -addr localhost:8080

go run cmd/nats-client/main.go -s localhost sb-events "ndjordjevic,new_order123"   

docker pull nats:latest   
docker run --name nats -p 4222:4222 --network event-hub-net -ti nats:latest

docker build -t ndjordjevic/server-echo -f cmd/server-echo/Dockerfile . (from the repo root)

docker run --name event-hub -tid --network event-hub-net -p 8080:8080 -e NATS_ADDR=localhost ndjordjevic/server-echo

docker service create --name event-hub --replicas 3 -e NATS_ADDR=192.168.99.1 -p 8080:8080 -l 'traefik.http.routers.event-hub.rule=Host(`event-hub.localhost`)' -l 'traefik.http.services.event-hub-service.loadbalancer.server.port=8080' ndjordjevic/server-echo

docker-machine scp event-hub.yml node1:/home/docker

nc -zvw3 192.168.99.100 4222

docker service logs event-hub_traefik -f

docker stack deploy -c event-hub.yml event-hub

docker service ls

docker stack rm event-hub

sudo netstat -tulpn | grep LISTEN


