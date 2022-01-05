# test-creatly-golang-dev

Task - [Task for golang dev (Creatly)](Task.md)

## Description

It is necessary to implement a web service that will have the following endpoints.

- POST /sign-up

Used for registration, accepts email and password at the entrance.

- POST /sign-in

Used for authentication, accepts an email and password at the entrance.   

- POST /upload

It is used to upload files that should later be uploaded to external Object Storage.

- GET /files

Returns information about all uploaded files (size, upload date, user ID, link to external storage).

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
