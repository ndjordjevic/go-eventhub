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

go tool pprof http://localhost:8080/debug/pprof/profile\?seconds\=30

kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.1/aio/deploy/recommended.yaml
kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | awk '/^deployment-controller-token-/{print $1}') | awk '$1=="token:"{print $2}'
kubectl proxy (start dashboard for docker desktop)
http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/#/login (url to docker desktop)

kubectl create -f deployments/nats-deployment.yaml
kubectl apply -f deployments/nats-deployment.yaml
kubectl logs nats-deployment-864cbddb96-stkqg

kubectl create -f deployments/nats-service.yaml
