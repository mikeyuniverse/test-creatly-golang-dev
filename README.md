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

or

```docker-compose
docker-compose up --build app
```

or

```make
make dcup
```

## Tests

|Package|Percent|
|---|---|
|TOTAL|83.2 %|
|Config|76.9 %|
|Handlers|86.4 %|
|Services|85.2 %|
