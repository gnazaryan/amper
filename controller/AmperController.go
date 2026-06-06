package controller

import (
	"amper/auth/authorization"
	"amper/common/structs"
	"amper/common/util"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// WidgetController is responsible for dispatching requests related to
// dashboard and underlying widgets managment functionalities
func AmperController(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	var resultStruct interface{}
	if len(pathSplit) > 2 {
		switch pathSplit[2] {
		case "fetch":
			resultStruct = fetchInstances(userID, w, r)
		case "remove":
			resultStruct = removeInstance(userID, w, r)
		case "create":
			resultStruct = createInstance(userID, sessionId, w, r)
		case "edit":
			resultStruct = editInstance(userID, sessionId, w, r)
		case "fetchInstanceInfo":
			resultStruct = fetchInstanceInfo(userID, sessionId, w, r)
		case "info":
			resultStruct = AmperInfo(userID)
		case "status":
			resultStruct = status(userID, w, r)
		case "invalidateCache":
			resultStruct = invalidateCache(userID, w, r)
		default:
		}
		marshaled, _ := json.Marshal(resultStruct)
		result = string(marshaled)
	}
	return
}

func fetchInstanceInfo(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result structs.AmperResult) {
	var amper structs.Amper
	json.NewDecoder(r.Body).Decode(&amper)
	info, err := authorization.FetchInstanceInfo(userID, sessionId, amper)
	if err == nil {
		result.Success = info != nil
		result.Data = info
	} else {
		result.Error = err.Error()
	}
	return result
}

func AmperInfo(userID *int64) (result structs.AmperResult) {
	var err error = nil
	if err == nil {
		info := structs.Amper{
			Id:      util.PointerInt64(1),
			Name:    util.PointerString("Amper"),
			Address: util.PointerString("localhost"),
		}
		result.Data = &info
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}

func fetchInstances(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.AmperResults) {
	var parameters struct {
		Type *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	instances, err := authorization.GetInstances(userID, parameters.Type)
	if err == nil {
		result.Data = instances
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}

func removeInstance(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var amper structs.Amper
	json.NewDecoder(r.Body).Decode(&amper)

	success, err := authorization.RemoveInstance(userID, amper)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return result
}

// create is responsible for creating a new user with the specified parametes
func createInstance(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var amper structs.Amper
	json.NewDecoder(r.Body).Decode(&amper)
	success, err := authorization.CreateInstance(userID, sessionId, amper)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return
}

func editInstance(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var amper structs.Amper
	json.NewDecoder(r.Body).Decode(&amper)
	success, err := authorization.EditInstance(userID, sessionId, amper)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return
}

func status(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.AmperResult) {
	var amper structs.Amper
	json.NewDecoder(r.Body).Decode(&amper)
	amperResult, errAR := authorization.FetchStatus(userID, amper)
	if errAR == nil {
		result.Success = true
		result.Data = amperResult
	} else {
		result.Success = false
		result.Error = errAR.Error()
	}
	return result
}

func invalidateCache(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Name          *string
		UserIdDelete  *string
		ChatChannelId *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	var UserIdDelete *int64
	if parameters.UserIdDelete != nil {
		UserIdDeleteInt, _ := strconv.ParseInt(*parameters.UserIdDelete, 10, 64)
		UserIdDelete = &UserIdDeleteInt
	}
	var ChatChannelDelete *int64
	if parameters.ChatChannelId != nil {
		ChatChannelDeleteInt, _ := strconv.ParseInt(*parameters.ChatChannelId, 10, 64)
		ChatChannelDelete = &ChatChannelDeleteInt
	}
	success, errIC := authorization.InvalidateCache(userID, parameters.Name, UserIdDelete, ChatChannelDelete)
	if errIC == nil {
		result.Success = success
	} else {
		result.Success = false
		result.Error = errIC.Error()
	}
	return result
}
