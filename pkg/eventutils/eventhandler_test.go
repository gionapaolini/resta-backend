package eventutils

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestHandleEvent(t *testing.T) {
	// Arrange
	settings, _ := esdb.ParseConnectionString(eventStoreConnectionString)
	db, _ := esdb.NewClient(settings)
	eventHandler := NewEventHandler(db, "IntegrationTestGroup")

	type testEvent1 struct {
		Data string
	}

	var dataReceived string

	var testHandler = func(rawEvent *esdb.SubscriptionEvent) {
		var event testEvent1
		json.Unmarshal(rawEvent.EventAppeared.Event.Data, &event)
		dataReceived = event.Data
	}

	data := testEvent1{
		Data: "some value 123",
	}

	//Act
	eventHandler.HandleEvent("testEvent1", testHandler).Start()
	sendTestEvent(db, "testStream", data)
	time.Sleep(1 * time.Second) // sleep to make sure the event is received

	// Assert
	require.Equal(t, data.Data, dataReceived)
}

func sendTestEvent(db *esdb.Client, streamName string, data interface{}) {

	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	options := esdb.AppendToStreamOptions{}

	db.AppendToStream(context.Background(), streamName, options, esdb.EventData{
		ContentType: esdb.JsonContentType,
		EventType:   utils.GetType(data),
		Data:        bytes,
	})
}
