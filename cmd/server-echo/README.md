go run -race main.go // start server

go run cmd/client-gorilla/main.go -user ndjordjevic

go run cmd/client-gorilla/main.go -user vpopovic

go run cmd/nats-client/main.go -s localhost sb-events "ndjordjevic,new_order123"   

docker pull nats:latest   

docker run -p 4222:4222 -ti nats:latest
