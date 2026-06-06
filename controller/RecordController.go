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
func RecordController(userID *int64, w *http.ResponseWriter, r *http.Request) (result string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	var resultStruct interface{}
	if len(pathSplit) > 2 {
		switch pathSplit[2] {
		case "fetch":
			resultStruct = fetchRecords(userID, w, r)
		case "add":
			resultStruct = addRecord(userID, w, r)
		case "remove":
			resultStruct = removeRecord(userID, w, r)
		case "update":
			resultStruct = updateRecord(userID, w, r)
		case "addRecords":
			resultStruct = addRecords(userID, w, r)
		case "removeRecords":
			resultStruct = removeRecords(userID, w, r)
		case "updateRecords":
			resultStruct = updateRecords(userID, w, r)
		default:
		}
		marshaled, _ := json.Marshal(resultStruct)
		result = string(marshaled)
	}
	return
}

func fetchRecords(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.RecordsResult) {
	var parameters struct {
		ApiName  *string
		ObjectId *int64
		Start    *int64
		Limit    *int64
		Search   *string
		Metadata *bool
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	records, metadata, totalCount, err := authorization.FetchRecords(userID, parameters.ApiName, parameters.ObjectId, parameters.Start, parameters.Limit, parameters.Search, parameters.Metadata)
	if err == nil {
		result.Data = records
		result.Metadata = metadata
		result.TotalCount = totalCount
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}

func addRecord(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.RecordResult) {
	var parameters struct {
		ApiName *string
		Payload *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	record, err := authorization.AddRecord(userID, parameters.ApiName, parameters.Payload)
	if err == nil {
		result.Data = record
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}

func removeRecord(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.RecordResult) {
	var parameters struct {
		Identifier *string
		Id         *int64
		ApiName    *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	record, err := authorization.RemoveRecord(userID, parameters.ApiName, parameters.Identifier, parameters.Id)

	if err == nil {
		result.Data = record
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}

func updateRecord(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.RecordResult) {
	var parameters struct {
		ApiName *string
		Payload *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	record, err := authorization.UpdateRecord(userID, parameters.ApiName, parameters.Payload)

	if err == nil {
		result.Data = record
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}

func addRecords(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.RecordsResult) {
	var parameters struct {
		ApiName *string
		Payload *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	recordSuccess, recordError, err := authorization.AddRecords(userID, parameters.ApiName, parameters.Payload)

	if err == nil {
		result.Data = recordSuccess
		result.ErrorData = recordError
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}

func removeRecords(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.RecordResult) {
	var parameters struct {
		Identifiers *[]string
		Ids         *[]int64
		ApiName     *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.RemoveRecords(userID, parameters.ApiName, parameters.Ids, parameters.Identifiers)

	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return result
}

func updateRecords(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.RecordsResult) {
	var parameters struct {
		Payloads *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	resultSuccess, resultError, err := authorization.UpdateRecords(userID, parameters.Payloads)

	if err == nil {
		result.Data = resultSuccess
		result.ErrorData = resultError
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return result
}
