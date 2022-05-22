package main

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/gofrs/uuid"
)

const eventStoreConnectionString = "esdb://127.0.0.1:2113?tls=false&keepAliveTimeout=10000&keepAliveInterval=10000"

func main() {
	settings, _ := esdb.ParseConnectionString(eventStoreConnectionString)
	db, _ := esdb.NewClient(settings)
	eventStore, _ := NewEventStore(db)

	normalMenuItem := NewMenu()
	preparedMenuItem := NewPreparedMenuItem()

	normalMenuItem.ChangeName("Jello1")
	preparedMenuItem.ChangeName("Jello2")
	preparedMenuItem.ChangePreparationTime(15 * time.Minute)

	streamName1 := "Menu-" + normalMenuItem.ID.String()
	streamName2 := "Menu-" + preparedMenuItem.ID.String()

	eventStore.SaveEventsToNewStream(streamName1, normalMenuItem.Events)
	eventStore.SaveEventsToNewStream(streamName2, preparedMenuItem.Events)

	retrievedEvents1, _ := eventStore.GetAllEventsByStreamName(streamName1)
	retrievedEvents2, _ := eventStore.GetAllEventsByStreamName(streamName2)

	menu1 := &MenuItem{}
	menu2 := &PreparedMenuItem{}
	reconstructFromEvents(menu1, retrievedEvents1)
	reconstructFromEvents(menu2, retrievedEvents2)
	print("hellooo")
}

type IEvent interface {
	GetEventID() uuid.UUID
}

type Event struct {
	ID   uuid.UUID
	Name string
	Data []byte
}

type Entity struct {
	ID        uuid.UUID
	Events    []Event
	NewEvents []Event
}

type MenuItem struct {
	Entity
	Name, ImageURL string
}

type IReconstructible interface {
	ApplyEvent(event Event)
}

func NewMenu() *MenuItem {
	menuCreatedEvent := MenuCreated{
		EventInfo: newEventInfo(getRandomID()),
		Name:      "Default",
		ImageURL:  "lollo.lol",
	}
	event := createSerializedEvent(menuCreatedEvent)
	menuItem := &MenuItem{}
	menuItem.ApplyEvent(event)
	menuItem.NewEvents = append(menuItem.NewEvents, event)
	return menuItem
}

func createSerializedEvent(eventObj IEvent) Event {
	bytes, _ := json.Marshal(eventObj)
	event := Event{
		ID:   eventObj.GetEventID(),
		Name: getType(eventObj),
		Data: bytes,
	}
	return event
}

func (mi *MenuItem) ChangeName(newName string) {
	menuNameChangedEvent := MenuNameChanged{
		EventInfo: newEventInfo(mi.ID),
		NewName:   newName,
	}
	event := createSerializedEvent(menuNameChangedEvent)
	mi.ApplyEvent(event)
	mi.NewEvents = append(mi.NewEvents, event)
}

func (mi *MenuItem) ApplyEvent(event Event) {
	switch event.Name {
	case "MenuCreated":
		var e MenuCreated
		json.Unmarshal(event.Data, &e)
		applyMenuCreated(mi, e)
	case "MenuNameChanged":
		var e MenuNameChanged
		json.Unmarshal(event.Data, &e)
		applyMenuNameChanged(mi, e)
	case "MenuImageURLChanged":
		var e MenuImageURLChanged
		json.Unmarshal(event.Data, &e)
		applyMenuImageURLChanged(mi, e)
	default:
		return
	}
	mi.Events = append(mi.Events, event)
}

func applyMenuCreated(mi *MenuItem, e MenuCreated) {
	mi.ID = e.EntityID
	mi.Name = e.Name
	mi.ImageURL = e.ImageURL
}

func applyMenuNameChanged(mi *MenuItem, e MenuNameChanged) {
	mi.Name = e.NewName
}

func applyMenuImageURLChanged(mi *MenuItem, e MenuImageURLChanged) {
	mi.ImageURL = e.NewImageURL
}

type PreparedMenuItem struct {
	*MenuItem
	EstimatedPreparationTime time.Duration
}

func NewPreparedMenuItem() *PreparedMenuItem {
	menuItem := NewMenu()

	menuItemSpecializedIntoPreparedMenuItem := MenuItemSpecializedIntoPreparedMenuItem{
		EventInfo:                newEventInfo(menuItem.ID),
		EstimatedPreparationTime: 10 * time.Minute,
	}
	event := createSerializedEvent(menuItemSpecializedIntoPreparedMenuItem)
	preparedMenuItem := &PreparedMenuItem{}
	preparedMenuItem.MenuItem = menuItem
	preparedMenuItem.ApplyEvent(event)
	preparedMenuItem.NewEvents = append(preparedMenuItem.NewEvents, event)
	return preparedMenuItem
}

func (mi *PreparedMenuItem) ChangePreparationTime(newTime time.Duration) {
	menuEstimatedPreparationTimeChanged := MenuEstimatedPreparationTimeChanged{
		EventInfo:                   newEventInfo(mi.ID),
		NewEstimatedPreparationTime: newTime,
	}
	event := createSerializedEvent(menuEstimatedPreparationTimeChanged)
	mi.ApplyEvent(event)
	mi.NewEvents = append(mi.NewEvents, event)
}

func (mi *PreparedMenuItem) ApplyEvent(event Event) {
	switch event.Name {
	case "MenuItemSpecializedIntoPreparedMenuItem":
		var e MenuItemSpecializedIntoPreparedMenuItem
		json.Unmarshal(event.Data, &e)
		applyMenuItemSpecializedIntoPreparedMenuItem(mi, e)
	case "MenuEstimatedPreparationTimeChanged":
		var e MenuEstimatedPreparationTimeChanged
		json.Unmarshal(event.Data, &e)
		applyMenuEstimatedPreparationTimeChanged(mi, e)
	default:
		if mi.MenuItem == nil {
			mi.MenuItem = &MenuItem{}
		}
		mi.MenuItem.ApplyEvent(event)
		return
	}
	mi.Events = append(mi.Events, event)
}

func applyMenuItemSpecializedIntoPreparedMenuItem(mi *PreparedMenuItem, e MenuItemSpecializedIntoPreparedMenuItem) {
	mi.EstimatedPreparationTime = e.EstimatedPreparationTime
}

func applyMenuEstimatedPreparationTimeChanged(mi *PreparedMenuItem, e MenuEstimatedPreparationTimeChanged) {
	mi.EstimatedPreparationTime = e.NewEstimatedPreparationTime
}

// Events

type EventInfo struct {
	EntityID, EventID uuid.UUID
	CreatedAt         time.Time
}

func (ei EventInfo) GetEventID() uuid.UUID {
	return ei.EventID
}

type MenuCreated struct {
	EventInfo
	Name, ImageURL string
}
type MenuItemSpecializedIntoPreparedMenuItem struct {
	EventInfo
	EstimatedPreparationTime time.Duration
}

type MenuNameChanged struct {
	EventInfo
	NewName string
}

type MenuImageURLChanged struct {
	EventInfo
	NewImageURL string
}

type MenuEstimatedPreparationTimeChanged struct {
	EventInfo
	NewEstimatedPreparationTime time.Duration
}

func getRandomID() uuid.UUID {
	id, _ := uuid.NewV4()
	return id
}

func reconstructFromEvents(entity IReconstructible, events []Event) {
	for _, event := range events {
		entity.ApplyEvent(event)
	}
}

func newEventInfo(entityID uuid.UUID) EventInfo {
	eventID := getRandomID()
	eventInfo := EventInfo{
		EventID:   eventID,
		EntityID:  entityID,
		CreatedAt: time.Now(), // change
	}
	return eventInfo
}

func getType(x interface{}) string {
	dataType := reflect.TypeOf(x).String()
	typeName := dataType[strings.LastIndex(dataType, ".")+1:]
	return typeName
}
