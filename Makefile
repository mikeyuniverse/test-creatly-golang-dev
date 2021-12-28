SILENT:

run:
	go run cmd/main.go

dcup:
	docker-compose up -d --remove-orphans app