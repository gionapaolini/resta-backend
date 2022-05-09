package main

// THIS IS THE INITIALIZATION SCRIPT FOR THE DEVELOPMENT ENVIRONMENT
// MAKE SURE EVERY ACTION IS IDEMPOTENT AND THAT IT PANICS ON ERRORS
// SO THAT IT CAN BE RESTARTED ON FAILURES (for example running before the DB is ready)

import (
	"context"
	"errors"
	"os"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

const eventStoreConnectionString = "esdb://127.0.0.1:2113?tls=false&keepAliveTimeout=10000&keepAliveInterval=10000"
const postgresConnectionString = "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"
const migrationsPath = "file:///src/service.menu/queries/internal/migrations"

// const migrationsPath = "file://../service.menu/queries/internal/migrations"

func main() {
	os.MkdirAll("/src/resources/images/categories", 0755)
	CreatePersistentSubscription("IntegrationTestGroup", []string{
		"testEvent1",
	})
	CreatePersistentSubscription("menu.queries", []string{
		"MenuCreated",
		"MenuEnabled",
		"MenuDisabled",
		"MenuNameChanged",
		"CategoryCreated",
		"CategoryAddedToMenu",
		"CategoryNameChanged",
		"SubCategoryCreated",
		"SubCategoryAddedToCategory",
	})
	CreatePersistentSubscription("menu.commands", []string{
		"CategoryCreated",
		"SubCategoryCreated",
	})
	RunPostgresMigrations()
}

func RunPostgresMigrations() {
	m, err := migrate.New(
		migrationsPath,
		postgresConnectionString,
	)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil {
		if err.Error() == "no change" {
			return
		}
	}
	panic(err)
}

func CreatePersistentSubscription(groupName string, events []string) {
	settings, _ := esdb.ParseConnectionString(eventStoreConnectionString)
	db, _ := esdb.NewClient(settings)

	err := db.CreatePersistentSubscriptionAll(
		context.Background(),
		groupName,
		esdb.PersistentAllSubscriptionOptions{
			Filter: &esdb.SubscriptionFilter{
				Type:     esdb.EventFilterType,
				Prefixes: events,
			},
		},
	)
	if err != nil {
		var badInputErr *esdb.PersistentSubscriptionError
		if errors.As(err, &badInputErr) {
			if badInputErr.Code == 6 {
				return
			}
			panic(badInputErr)
		}
		panic(err)
	}
}
