package controller

import (
	"amper/auth/authorization"
	"amper/common/structs"
	"encoding/json"
	"net/http"
	"strings"
)

// UserController is responsible for dispatching requests related to
// user managment functionalities
func EntityTypeController(userID *int64, w *http.ResponseWriter, r *http.Request) (result string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	var resultStruct interface{}
	if len(pathSplit) > 2 {
		switch pathSplit[2] {
		case "getObjectTypes":
			resultStruct = getObjectTypes(userID, w, r)
		case "getObjectTypeFields":
			resultStruct = getObjectTypeFields(userID, w, r)
		case "createObjectType":
			resultStruct = createObjectTypes(userID, w, r)
		case "deleteObjectType":
			resultStruct = deleteObjectType(userID, w, r)
		default:
		}
		marshaled, _ := json.Marshal(resultStruct)
		result = string(marshaled)
	}
	return
}

func deleteObjectType(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		ObjectTypeId *int64
	}

	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.DeleteObjectType(userID, parameters.ObjectTypeId)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return
}

func createObjectTypes(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		ObjectId  *int64
		ExtendsTo *int64
		ApiName   *string
		Label     *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.CreateObjectType(userID, parameters.ObjectId, parameters.ExtendsTo, parameters.ApiName, parameters.Label)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return
}

// getObjectTypes is resposible for retrieving object types for the provided object type id
// user requesting the information should have permissions
func getObjectTypes(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ObjectTypes) {
	var parameters struct {
		ObjectId *int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	objectTypes, err := authorization.GetObjectTypes(userID, parameters.ObjectId)
	if err == nil {
		result.Data = objectTypes
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return
}

// getObjectTypes is resposible for retrieving object types for the provided object type id
// user requesting the information should have permissions
func getObjectTypeFields(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ObjectTypes) {
	var parameters struct {
		ObjectId     *int64
		ObjectTypeId *int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	objectType, err := authorization.GetObjectTypeFields(userID, parameters.ObjectId, parameters.ObjectTypeId)
	if err == nil {
		Data := make([]structs.ObjectType, 1)
		Data[0] = *objectType
		result.Data = &Data
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return
}
