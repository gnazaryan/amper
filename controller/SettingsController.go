package controller

import (
	"amper/auth/authorization"
	"amper/common/structs"
	"encoding/json"
	"net/http"
	"strings"
)

func SettingsController(userID *int64, w *http.ResponseWriter, r *http.Request) (result string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	var resultStruct interface{}
	if len(pathSplit) > 2 {
		switch pathSplit[2] {
		case "save":
			resultStruct = saveSettings(userID, w, r)
		case "fetch":
			resultStruct = fetchSettings(userID, w, r)
		default:
		}
		marshaled, _ := json.Marshal(resultStruct)
		result = string(marshaled)
	}
	return
}

func saveSettings(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Settings *structs.Settings `json:"settings"`
	}
	json.NewDecoder(r.Body).Decode(&parameters)

	success, err := authorization.SaveSettings(userID, parameters.Settings)
	if err == nil && success {
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}

func fetchSettings(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.SettingsResult) {
	settings, err := authorization.FetchSettings(userID)
	if err == nil {
		result.Success = true
		result.Data = settings
	} else {
		result.Error = err.Error()
	}
	return result
}
