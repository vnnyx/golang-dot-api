# GOLANG-DOT-API

## Tech stack

- Golang v1.19
- Echo V4
- Redis for caching
- Golang ozzo-validator
- Google wire for dependency injector
- Mockery for mock generator
- Docker and docker compose for deployment
- Jenkins for deployment

## Design Pattern

**Model > Repository > Service > Controller (Echo REST API)**

The design pattern above is an implementation of clean architecture by uncle bob. The clean architecture has several layers: entities, use cases, controllers, and frameworks & drivers. The reason for using this design pattern is that it is easy to maintain.

## Setup

1. copy .env.example to .env
2. run go mod tidy
3. run main app with `go run cmd/app/main.go`

## Live Demo

I deployed this service, and you can access it via `https://cloud.vnnyx.my.id/dot-api/{ENDPOINT}`

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

## Testing

### Integration Test

![integration_test](https://imgur.com/0EdCmp3.png)

### Unit Test

![unit_test](https://imgur.com/05MuWmO.png)

## Deployment

![deployment](https://imgur.com/oliaMaR.png)

## API Docs
[API Docs](https://documenter.getpostman.com/view/24450154/2s8YzWRfeC)
