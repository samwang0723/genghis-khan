# genghis-khan
Chatbot purchasing service based on honestbee APIs

### Docker-compose

    Three step setup & running project:
    1. make build # Build base docker image for running docker-compose
    2. make up # Execute `docker-compose up`
    3. make start # Execute `go run main.go` in docker-compose

    One step stop the project:
    1. make down # Stop and clear docker env