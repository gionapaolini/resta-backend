# Resta - Backend webservices

To start the development environment (Postgres, EventStore):

```
cd dev
sudo docker-compose up --build -d 
```

To start the menu.queries service

```
cd service.menu\queries
go run main.go
```

To start the menu.commands service

```
cd service.menu\commands
go run main.go
```