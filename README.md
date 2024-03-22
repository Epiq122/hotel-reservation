# Hotel reservation backend

## Project outline

- users -> book room from an hotel
- admins -> check reservations and bookings
- authentication and authorization -> JWT tokens
- hotels -> CRUD API -> JSON
- rooms -> CRUD API -> JSON
- scripts -> database management -> seeding, migration

## Resources

### MongoDB driver

Documentation

```
https://mongodb.com/docs/drivers/go/current/quick-start
```

Installing mongoDB client

```
go get go.mongodb.org/mongo-driver/mongo
```

## gofiber

```
https://gofiber.io
```

Installing gofiber

```
go get -u github.com/gofiber/fiber/v2
```

## Docker

## Installing mongodb as a docker container

```
docker run -name mongodb -d mongo:latest -p 27017:27017
```
