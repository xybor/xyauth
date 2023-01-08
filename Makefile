DATABASE=./database
SERVER=./server
TEMPLATE=./template

DOCKER=docker
DOCKER_COMPOSE=docker compose
DOCKERFILE=build/Dockerfile
DOCKER_COMPOSE_FILE=deploy/compose.yaml

IMAGE_NAME=xyauth

.PHONY: build run clean
.PHONY: docker-gen docker-build docker-start docker-stop docker-clean
.PHONY: cert-gen

database:
	go build -o $(DATABASE) ./cmd/database/*.go

template:
	go build -o $(TEMPLATE) ./cmd/template/*.go

server:
	go build -o $(SERVER) ./cmd/server/*.go

clean:
	rm -f $(DATABASE) $(SERVER) $(TEMPLATE)

run: server
	$(SERVER)

docker-gen: template
	$(TEMPLATE) $(DOCKER_COMPOSE_FILE).template -c configs/compose.ini

docker-build:
	$(DOCKER) build -t $(IMAGE_NAME) -f $(DOCKERFILE) .

docker-start:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d

docker-stop:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down

docker-access:
	$(DOCKER) exec -it ${IMAGE_NAME} sh

docker-clean: docker-stop
	rm -f $(DOCKER_COMPOSE_FILE)
	$(DOCKER) rmi -f $(IMAGE_NAME)

cert-gen:
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.crt -subj \
		"/C=VN/ST=HoChiMinh/L=GV/O=xybor/OU=XyborSpace/CN=localhost"
	@echo "The output is server.key and server.crt"

test:
	go test -race ./...
