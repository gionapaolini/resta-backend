package utils

import (
	"encoding/json"
	"reflect"
	"strings"
)

type TypedJson struct {
	Type string
	Data json.RawMessage
}

func GetType(x interface{}) string {
	dataType := reflect.TypeOf(x).String()
	typeName := dataType[strings.LastIndex(dataType, ".")+1:]
	return typeName
}

func SerializeObject(obj interface{}) []byte {
	data, err := json.Marshal(obj)
	typeName := GetType(obj)
	typedThing := TypedJson{typeName, data}
	finalJson, _ := json.Marshal(typedThing)
	if err != nil {
		panic(err)
	}
	return finalJson
}
