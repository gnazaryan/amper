package controller

import (
	"amper/auth/authorization"
	"amper/common/structs"
	"encoding/json"
	"net/http"
	"strings"
)

func UpdatesController(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	var resultStruct interface{}
	if len(pathSplit) > 2 {
		switch pathSplit[2] {
		case "fetch":
			resultStruct = fetchUpdates(userID, sessionId, w, r)
		case "push":
			resultStruct = pushUpdates(userID, w, r)
		default:
		}
		marshaled, _ := json.Marshal(resultStruct)
		result = string(marshaled)
	}
	return
}

func fetchUpdates(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result structs.UserUpdatesResult) {
	userUpdates, errUU := authorization.FetchUpdates(userID)
	if errUU == nil {
		result.Success = true
		result.Data = userUpdates
	} else {
		result.Error = errUU.Error()
	}
	return result
}

func pushUpdates(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Category     *string
		Participants *[]int64
		Value        *interface{}
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, errPU := authorization.PutUpdates(userID, parameters.Category, parameters.Participants, parameters.Value)
	if errPU == nil && success {
		result.Success = true
	} else {
		result.Error = errPU.Error()
	}
	return result
}
