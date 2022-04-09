Run the eventstore:
```
sudo docker run --rm -d  --name esdb-node -it -p 2113:2113 -p 1113:1113 eventstore/eventstore:latest --insecure --run-projections=All --enable-atom-pub-over-http
```

Run postgres:
```
sudo docker run --rm -d --name postgres14 -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword postgres:14
```

Create a migration:

```
migrate create -ext sql -dir {migrationDIR} -seq {migrationName}
```

Run migrations for views:
```
migrate -database postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable -path internal/migrations up
```