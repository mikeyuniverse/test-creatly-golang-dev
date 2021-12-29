SILENT:

run:
	go run cmd/main.go

dcup:
	dbuild
	docker-compose up -d --remove-orphans app

build:
	go build -o ./bin/app ./cmd/main.go

dbuild:
	docker build . -t=creatly-dev

drun:
	docker run creatly-dev

ddel:
	docker image rm -f creatly-dev

dreload:
	make ddel
	make dbuild
	make drun