package controller

import (
	"amper/auth/authorization"
	"amper/common/structs"
	"encoding/json"
	"net/http"
	"strings"
)

func FieldController(userID *int64, w *http.ResponseWriter, r *http.Request) (result string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	var resultStruct interface{}
	if len(pathSplit) > 2 {
		switch pathSplit[2] {
		case "createField":
			resultStruct = createField(userID, w, r)
		case "deleteField":
			resultStruct = deleteField(userID, w, r)
		case "addObjectTypeField":
			resultStruct = addObjectTypeField(userID, w, r)
		case "deleteObjectTypeField":
			resultStruct = deleteObjectTypeField(userID, w, r)
		}
		marshaled, _ := json.Marshal(resultStruct)
		result = string(marshaled)
	}
	return
}

func deleteObjectTypeField(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		FieldId      *int64
		ObjectTypeId *int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.DeleteObjectTypeField(userID, parameters.FieldId, parameters.ObjectTypeId)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return
}

func addObjectTypeField(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		FieldId      *int64
		ObjectTypeId *int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.AddObjectTypeField(userID, parameters.FieldId, parameters.ObjectTypeId)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return
}

func deleteField(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		EntityId *int64
		FieldIds *[]int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.DeleteField(userID, parameters.EntityId, parameters.FieldIds)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return
}

func createField(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		EntityId        *int64
		ApiName         *string
		Label           *string
		Required        *bool
		Status          *bool
		DataType        *string
		TextLength      *int64
		ObjectReference *int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)

	success, err := authorization.CreateField(userID, parameters.EntityId, parameters.ApiName, parameters.Label, parameters.DataType, parameters.TextLength, parameters.ObjectReference, parameters.Required, parameters.Status)

	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return
}
