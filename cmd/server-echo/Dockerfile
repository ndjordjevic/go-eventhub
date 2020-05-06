FROM golang:stretch AS build-env
WORKDIR /go/src/github.com/ndjordjevic/go-eventhub/
COPY . ./
RUN ls
RUN CGO_ENABLED=0 GOOS=linux go build -o cmd/server-echo/server-echo github.com/ndjordjevic/go-eventhub/cmd/server-echo
RUN ls

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/ndjordjevic/go-eventhub/cmd/server-echo/server-echo /app/
ENTRYPOINT ./server-echo
