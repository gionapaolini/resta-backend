package eventutils

import (
	"encoding/json"

	"github.com/Resta-Inc/resta/pkg/utils"
)

type TestEntity struct {
	Entity
	State TestState
}

type TestState struct {
	Name string
}

func NewTestEntity() *TestEntity {
	var event IEvent
	entityID := utils.GenerateNewUUID()

	event = TestEntityCreated{
		EventInfo: NewEventInfo(entityID),
		Name:      "TestEvent",
	}

	testEntity := &TestEntity{}
	testEntity.SetNew()
	AddEvent(event, testEntity)
	return testEntity
}

func (testEntity *TestEntity) ChangeName(name string) {
	newEvent := TestEntityNameChanged{
		EventInfo: NewEventInfo(testEntity.ID),
		NewName:   name,
	}
	AddEvent(newEvent, testEntity)
}

type TestEntityCreated struct {
	EventInfo
	Name string
}

type TestEntityNameChanged struct {
	EventInfo
	NewName string
}

func (testEntity *TestEntity) ApplyEvent(event IEvent) {
	eventType := utils.GetType(event)
	switch eventType {
	case "TestEntityCreated":
		testEntity.applyTestEntityCreated(event.(TestEntityCreated))
	case "TestEntityNameChanged":
		testEntity.applyTestEntityNameChanged(event.(TestEntityNameChanged))
	}
}

func (testEntity *TestEntity) applyTestEntityCreated(event TestEntityCreated) {
	testEntity.State.Name = event.Name
	testEntity.ID = event.EntityID
}

func (testEntity TestEntity) applyTestEntityNameChanged(event TestEntityNameChanged) {
	testEntity.State.Name = event.NewName
}

func (testEntity TestEntity) DeserializeEvent(event Event) IEvent {
	switch event.Name {
	case "TestEntityCreated":
		var e TestEntityCreated
		json.Unmarshal(event.Data, &e)
		return e
	case "TestEntityNameChanged":
		var e TestEntityNameChanged
		json.Unmarshal(event.Data, &e)
		return e
	default:
		return nil
	}
}
