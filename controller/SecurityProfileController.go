package controller

import (
	"amper/auth/authorization"
	"amper/common/structs"
	"encoding/json"
	"net/http"
	"strings"
)

// WidgetController is responsible for dispatching requests related to
// dashboard and underlying widgets managment functionalities
func SecurityProfileController(userID *int64, w *http.ResponseWriter, r *http.Request) (result string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	var resultStruct interface{}
	if len(pathSplit) > 2 {
		switch pathSplit[2] {
		case "fetch":
			resultStruct = fetchProfiles(userID, w, r)
		default:
		}
		marshaled, _ := json.Marshal(resultStruct)
		result = string(marshaled)
	}
	return
}

func fetchProfiles(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ProfileResult) {
	var parameters struct {
		Start         *int
		Limit         *int
		Search        *[]string
		SorfField     *string
		SortDirection *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	userProfiles, err := authorization.FetchProfiles(userID, parameters.Start, parameters.Limit, parameters.Search, parameters.SorfField, parameters.SortDirection)
	if err == nil {
		result.Data = userProfiles
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}
