FROM golang:1.17-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

# build go app
RUN go mod download
RUN go build -o ./bin/app ./cmd/main.go

CMD ["./bin/app"]