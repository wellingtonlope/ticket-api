# Ticket API
This project aims to learn and put DDD and Clean Architecture into practice. 
With this in mind, I built a simple project where a user with a question puts his problem and an operator can visualize it and propose a solution.

Demo: https://wellingtonlope-ticket-api.herokuapp.com \
Documentation: https://wellingtonlope-ticket-api.herokuapp.com/swagger/index.html

![Build](https://github.com/wellingtonlope/ticket-api/actions/workflows/build.yaml/badge.svg)
![Coverage](https://img.shields.io/badge/Coverage-32.4%25-yellow)

## Usage
You can use .env to set up the environment and start the http server with the following command:
```bash
go run cmd/main.go
```

Also, you can start the server with docker and docker-compose, just run the following command:
```bash
# by default the server will listen on port 1323
docker-compose up --build
```

## Tests
To run the tests, run the following command:
```bash
go test ./...
```