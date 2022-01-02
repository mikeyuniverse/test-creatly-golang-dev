SILENT:

run:
	go run cmd/main.go



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

dcup:
	make dbuild
	docker-compose up -d --remove-orphans app

test:
	go test -coverprofile=cover.out ./...
	go tool cover -func=cover.out
