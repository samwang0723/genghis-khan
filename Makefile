OUT := bin/genghis-khan

all: install

install: deps
		dep ensure -v

.PHONY: test
test:
	CGO_ENABLED=0 go test -v github.com/samwang0723/genghis-khan/honestbee

deps:
		@hash dep > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
			go get github.com/golang/dep/cmd/dep; \
		fi

.PHONY: build
build:
	docker-compose build

.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.PHONY: start
start:
	docker exec -it genghis-khan_api_1 go run main.go