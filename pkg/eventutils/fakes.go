package eventutils

import (
	"encoding/json"

	"github.com/Resta-Inc/resta/pkg/utils"
)

type TestEntity struct {
	*Entity
	State TestState
}

type TestState struct {
	Name string
}

func NewTestEntity() TestEntity {
	var event IEvent
	entityID := utils.GenerateNewUUID()

	testEntity := EmptyTestEntity()

	event = TestEntityCreated{
		EntityEventInfo: NewEntityEventInfo(entityID),
		Name:            "TestEvent",
	}
	return AddNewEvent(testEntity, event).(TestEntity)
}

func EmptyTestEntity() TestEntity {
	return TestEntity{
		Entity: &Entity{},
	}
}

func (testEntity TestEntity) ChangeName(name string) TestEntity {
	newEvent := TestEntityNameChanged{
		EntityEventInfo: NewEntityEventInfo(testEntity.ID),
		NewName:         name,
	}
	testEntity = AddNewEvent(testEntity, newEvent).(TestEntity)

	return testEntity
}

type TestEntityCreated struct {
	EntityEventInfo
	Name string
}

type TestEntityNameChanged struct {
	EntityEventInfo
	NewName string
}

func (testEntity TestEntity) ApplyEvent(event IEntityEvent) IEntity {
	eventType := utils.GetType(event)
	switch eventType {
	case "TestEntityCreated":
		testEntity = testEntity.applyTestEntityCreated(event.(TestEntityCreated))
	case "TestEntityNameChanged":
		testEntity = testEntity.applyTestEntityNameChanged(event.(TestEntityNameChanged))
	}
	return testEntity
}

func (testEntity TestEntity) applyTestEntityCreated(event TestEntityCreated) TestEntity {
	testEntity.State.Name = event.Name
	testEntity.ID = event.EntityID
	return testEntity
}

func (testEntity TestEntity) applyTestEntityNameChanged(event TestEntityNameChanged) TestEntity {
	testEntity.State.Name = event.NewName
	return testEntity
}

func (testEntity TestEntity) DeserializeEvent(jsonData []byte) IEvent {
	eventType, rawData := GetRawDataFromSerializedEvent(jsonData)
	switch eventType {
	case "TestEntityCreated":
		var e TestEntityCreated
		json.Unmarshal(rawData, &e)
		return e
	case "TestEntityNameChanged":
		var e TestEntityNameChanged
		json.Unmarshal(rawData, &e)
		return e
	default:
		return nil
	}
}
