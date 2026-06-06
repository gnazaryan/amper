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
func WidgetController(userID *int64, w *http.ResponseWriter, r *http.Request) (result string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	var resultStruct interface{}
	if len(pathSplit) > 2 {
		switch pathSplit[2] {
		case "fetch":
			resultStruct = fetchWidgets(userID, w, r)
		case "add":
			resultStruct = addWidget(userID, w, r)
		case "remove":
			resultStruct = removeWidget(userID, w, r)
		case "update":
			resultStruct = updateWidget(userID, w, r)
		case "interactions":
			resultStruct = iteractions(userID, w, r)
		default:
		}
		marshaled, _ := json.Marshal(resultStruct)
		result = string(marshaled)
	}
	return
}

func iteractions(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.DashboardResult) {
	var parameters struct {
		DashboardId   *int64
		WidgetId      *int64
		ObjectApiName *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	dashboardWidgets, err := authorization.GetInteractions(userID, parameters.DashboardId, parameters.WidgetId, parameters.ObjectApiName)
	if err == nil {
		result.Data = dashboardWidgets
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}

func fetchWidgets(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.DashboardResult) {
	var parameters struct {
		DashboardId *int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	dashboardWidgets, err := authorization.GetWidgets(userID, parameters.DashboardId)
	if err == nil {
		result.Data = dashboardWidgets
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}

func addWidget(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.DashboardResult) {
	var parameters struct {
		DashboardId   *int64
		Label         *string
		Description   *string
		Configuration *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.AddWidget(userID, parameters.DashboardId, parameters.Label, parameters.Description, parameters.Configuration)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return result
}

func removeWidget(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.DashboardResult) {
	var parameters struct {
		DashboardId *int64
		WidgetId    *int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.RemoveWidget(userID, parameters.DashboardId, parameters.WidgetId)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return result
}

func updateWidget(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.DashboardResult) {
	var parameters structs.Dashboard

	json.NewDecoder(r.Body).Decode(&parameters)

	success, err := authorization.UpdateWidget(userID, &parameters)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return result
}
