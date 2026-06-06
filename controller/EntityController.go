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
func EntityController(userID *int64, w *http.ResponseWriter, r *http.Request) (result string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	var resultStruct interface{}
	if len(pathSplit) > 2 {
		switch pathSplit[2] {
		case "getEntities":
			resultStruct = getEntities(userID, w, r)
		case "getEntity":
			resultStruct = getEntity(userID, w, r)
		case "deleteEntity":
			resultStruct = deleteEntity(userID, w, r)
		case "create":
			resultStruct = createEntity(userID, w, r)
		case "edit":
			resultStruct = editEntity(userID, w, r)
		case "getFields":
			resultStruct = getFields(userID, w, r)
		default:
		}
		marshaled, _ := json.Marshal(resultStruct)
		result = string(marshaled)
	}
	return
}

// getFields is resposible for retrieving entity fields with the provide entity id parameter and the
// user requesting the information should have permissions
func getFields(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.FieldsResult) {
	var parameters struct {
		ObjectId *int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	fields, err := authorization.GetFields(userID, parameters.ObjectId)
	if err == nil {
		result.Data = fields
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return
}

// getEntities is resposible for retrieving entities with the provide start and limit parameters and the
// user requesting the information should have permissions
func getEntities(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Entities) {
	var parameters struct {
		Start *int
		Limit *int
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	entities, err := authorization.GetEntities(userID, parameters.Start, parameters.Limit)
	if err == nil {
		result.Data = entities
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return
}

// getEntity is resposible for retrieving the entity with the provided id
// user requesting the information should have permissions
func getEntity(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.EntityResponse) {
	var parameters struct {
		EntityId *int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	entity, err := authorization.GetEntity(userID, parameters.EntityId)
	if err == nil {
		result.Entity = entity
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return
}

// deleteEntity is resposible for deleting entity with the provide id
// user requesting the information should have permissions
func deleteEntity(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		EntityId *int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.DeleteEntity(userID, parameters.EntityId)

	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()

	}
	return
}

// create is resposible for deleting entity with the provided information
// user requesting the information should have permissions
func createEntity(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		EntityId    *int64
		ApiName     *string
		Title       *string
		TitlePlural *string
	}

	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.CreateEntity(userID, parameters.ApiName, parameters.Title, parameters.TitlePlural)

	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return
}

// edit is resposible for editing entity with the provided information
// user requesting the information should have permissions
func editEntity(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Id          *int64
		ApiName     *string
		Title       *string
		TitlePlural *string
	}

	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.EditEntity(userID, parameters.Id, parameters.ApiName, parameters.Title, parameters.TitlePlural)

	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return
}
