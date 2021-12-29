# test-creatly-golang-dev

Task - [Task for golang dev (Creatly)](Task.md)

## Description

It is necessary to implement a web service that will have the following endpoints.

- POST /sign-up
- POST /sign-in
- POST /upload
- GET /files

## Run

```go
go run cmd/main.go
```

или

```docker-compose
docker-compose up --build app
```

или

```make
make dcup
```

## Tests

Comming soon

<!-- 
## Commands

Build docker image

```dockerfile
docker build -t creatly-dev .
```

Run docker image

```docker
docker run -d -p 8000:8000 creatly-dev
```

Run docker-compose

```docker-compose
docker-compose up -d --remove-orphans app
``` -->