package MUtils

import (
	"reflect"
	"strconv"
	"fmt"
	"encoding/json"
	"database/sql/driver"
)

func IdPanic(id interface{}) string {
	typeInt64 := reflect.TypeOf((int64)(0))

	var idS string
	if reflect.TypeOf(id).Kind() == reflect.String {
		idS = id.(string)
	} else if reflect.TypeOf(id).ConvertibleTo(typeInt64) {	//float int int64 etc.
		idS = strconv.FormatInt(reflect.ValueOf(id).Convert(typeInt64).Int(), 10)
	} else {
		panic(fmt.Errorf("ID[%s] have illegal type", id))
	}

	return idS
}

func BytesPanic(mp map[string]interface{}) []byte {
	bs, err := json.Marshal(mp)
	if err != nil {
		panic(err)
	}
	return bs
}

type CanNull interface {
	Value() (driver.Value, error)
}

func CanNullToInterface(cn CanNull) interface{} {
	i, _ := cn.Value()
	return (interface{})(i)
}