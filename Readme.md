# Resta - Backend webservices

## To start the development environment:

```
cd dev
sudo docker-compose up --build -d 
```

This will run the EventStore and the Postgres server. 
Then it will run a script that will run the migrations for postgres and it will create all the needed subscriptions for the event store.
Check ./dev/main.go for details on the process

## To start the menu.queries service

```
cd service.menu\queries
go run main.go
```

## To start the menu.commands service

```
cd service.menu\commands
go run main.go
```