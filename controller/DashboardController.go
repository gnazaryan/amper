package controller

import (
	"amper/auth/authorization"
	"amper/common/structs"
	"encoding/json"
	"net/http"
	"strings"
)

// DashboardController is responsible for dispatching requests related to
// dashboard managment functionalities
func DashboardController(userID *int64, w *http.ResponseWriter, r *http.Request) (result string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	var resultStruct interface{}
	if len(pathSplit) > 2 {
		switch pathSplit[2] {
		case "fetch":
			resultStruct = fetchDashboards(userID, w, r)
		case "add":
			resultStruct = addDashboard(userID, w, r)
		case "remove":
			resultStruct = removeDashboard(userID, w, r)
		case "update":
			resultStruct = updateDashboard(userID, w, r)
		default:
		}
		marshaled, _ := json.Marshal(resultStruct)
		result = string(marshaled)
	}
	return
}

func fetchDashboards(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.DashboardsResult) {
	dashboards, err := authorization.GetDashboards(userID)
	if err == nil {
		result.Data = dashboards
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}

func addDashboard(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Label         *string
		Description   *string
		Configuration *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.AddDashboard(userID, parameters.Label, parameters.Description, parameters.Configuration)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()

	}
	return result
}

func updateDashboard(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Id 			  *int64
		Label         *string
		Description   *string
		Configuration *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.UpdateDashboard(userID, parameters.Id, parameters.Label, parameters.Description, parameters.Configuration)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()

	}
	return result
}

func removeDashboard(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		DashboardId *int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.RemoveDashboard(userID, parameters.DashboardId)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()

	}
	return result
}
