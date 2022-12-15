# GOLANG-DOT-API

## Tech stack

- Golang v1.19
- Echo V4
- Redis for caching
- Golang ozzo-validator
- gogle wire for dependency injector
- Mockery for mock generator
- Docker and docker compose for deployment
- Jenkins for deployment

## Design Pattern

**Model > Repository > Service > Controller (Echo REST API)**

Design pattern di atas merupakan implementasi dari clean architecture by uncle bob. Clean architecture sendiri memiliki beberapa layer yaitu entities, use case, controller, dan framework & driver. Alasan menggunakan design pattern tersebut adalah mudah dalam maintenance nya.

## Setup

1. copy .env.example to .env
2. go mod tidy
3. run main app with `go run cmd/app/main.go`

## Live Demo

Saya melakukan deployment terhadap service ini dan bisa diakses melalui `https://cloud.vnnyx.my.id/dot-api/{ENDPOINT}`

## Data For Login

```
username: user_sample
password: password
```

## List Endpoint

```
BASEURL: https://cloud.vnnyx.my.id/dot-api

POST /login
POST /logout

POST /user
GET /user/:id
GET /user
PUT /user/:id
DELETE /user/:id

POST /transaction
GET /transaction/id
GET /transaction
GET /transaction/user
PATCH /transaction/id
DELETE /transaction/id

```

## Tesing

### Integration Test

![integration_test](https://imgur.com/0EdCmp3.png)

## API Docs
[API Docs](https://documenter.getpostman.com/view/24450154/2s8YzWRfeC)