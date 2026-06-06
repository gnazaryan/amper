package jsons

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

func IsJsonArray(data *string) bool {
	payloadBytes := []byte(*data)
	trimedData := bytes.TrimLeft(payloadBytes, " \t\r\n")
	isArray := len(trimedData) > 0 && trimedData[0] == '['
	return isArray
}

func IsJsonObject(data *string) bool {
	payloadBytes := []byte(*data)
	trimedData := bytes.TrimLeft(payloadBytes, " \t\r\n")
	isObject := len(trimedData) > 0 && trimedData[0] == '{'
	return isObject
}

func GetJsonObject(data *string) (result map[string]interface{}, err error) {
	if !IsJsonObject(data) {
		err = fmt.Errorf("ths supplied data is not of a json object type")
		return
	}
	jsonDataReader := strings.NewReader(*data)
	decoder := json.NewDecoder(jsonDataReader)
	errJ := decoder.Decode(&result)
	if errJ != nil {
		err = fmt.Errorf("the supplied data is of wrong json object structure")
		return
	}
	return
}

func GetJsonArray(data *string) (result []map[string]interface{}, err error) {
	if !IsJsonArray(data) {
		err = fmt.Errorf("ths supplied data is not of a json array type")
		return
	}
	jsonDataReader := strings.NewReader(*data)
	decoder := json.NewDecoder(jsonDataReader)
	errJ := decoder.Decode(&result)
	if errJ != nil {
		err = fmt.Errorf("the supplied data is of wrong json array structure")
		return
	}
	return
}
