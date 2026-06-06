package controller

import (
	"amper/auth/authorization"
	"amper/common/structs"
	"amper/common/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// WidgetController is responsible for dispatching requests related to
// dashboard and underlying widgets managment functionalities
func ProfileController(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	var resultStruct interface{}
	if len(pathSplit) > 2 {
		switch pathSplit[2] {
		case "updateCover":
			resultStruct = updateCover(userID, w, r)
		case "viewCover":
			viewCover(userID, w, r)
		case "updatePhoto":
			resultStruct = updatePhoto(userID, w, r)
		case "viewPhoto":
			viewPhoto(userID, w, r)
		case "adjustPhoto":
			resultStruct = adjustPhoto(userID, sessionId, w, r)
		case "state":
			resultStruct = profileState(userID, w, r)
		case "saveConfiguration":
			resultStruct = saveConfiguration(userID, w, r)
		case "saveDetail":
			resultStruct = saveDetail(userID, w, r)
		case "addRelationship":
			resultStruct = addRelationship(userID, w, r)
		case "removeRelationship":
			resultStruct = removeRelationship(userID, w, r)
		default:
		}
		marshaled, _ := json.Marshal(resultStruct)
		result = string(marshaled)
	}
	return
}

func addRelationship(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		EmployeeId int64
		Type       *string
		Value      int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.AddRelationship(userID, parameters.EmployeeId, parameters.Type, parameters.Value)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return result
}

func removeRelationship(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		EmployeeId int64
		Type       *string
		Value      int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.RemoveRelationship(userID, parameters.EmployeeId, parameters.Type, parameters.Value)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return result
}

func saveDetail(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Name  *string
		Value *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.SaveDetail(userID, parameters.Name, parameters.Value)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return result
}

func saveConfiguration(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Name  *string
		Value *map[string]interface{}
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.SaveConfiguration(userID, parameters.Name, parameters.Value)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return result
}

func profileState(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ProfileConfigurationResult) {
	data, err := authorization.GetProfileState(userID)
	if err == nil {
		result.Success = true
		result.Data = data
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return result
}

func adjustPhoto(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result structs.ResultValue) {
	var parameters struct {
		PositionX *int
		PositionY *int
		Width     *int
		Height    *int
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, croppedImage, err := authorization.AdjustPhoto(userID, sessionId, parameters.PositionX, parameters.PositionY, parameters.Width, parameters.Height)
	if err == nil && success {
		result.Success = success
		result.Value = croppedImage
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return result
}

func viewCover(userID *int64, w *http.ResponseWriter, r *http.Request) {
	reader, metadata, err := authorization.ViewCover(userID)

	if err == nil && reader != nil {
		(*w).Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", *metadata.Name))
		(*w).Header().Set("Content-Type", util.IfElse(metadata.Rendition, *metadata.RenditionType, *metadata.Type).(string))
		io.Copy((*w), *reader)
	}
}

func updateCover(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ProfileResult) {
	data, errBR := ioutil.ReadAll(r.Body)
	if errBR != nil {
		result.Error = "we prefer to receive file content with the update cover request on the request body"
		result.Success = false
		return
	}
	success, err := authorization.UpdateCover(userID, data)
	if err == nil && success {
		result.Success = success
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return result
}

func viewPhoto(userID *int64, w *http.ResponseWriter, r *http.Request) {
	reader, metadata, err := authorization.ViewPhoto(userID)

	if err == nil && reader != nil {
		(*w).Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", *metadata.Name))
		(*w).Header().Set("Content-Type", util.IfElse(metadata.Rendition, *metadata.RenditionType, *metadata.Type).(string))
		io.Copy((*w), *reader)
	}
}

func updatePhoto(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ProfileResult) {
	data, errBR := ioutil.ReadAll(r.Body)
	if errBR != nil {
		result.Error = "we prefer to receive file content with the update cover request on the request body"
		result.Success = false
		return
	}
	success, err := authorization.UpdatePhoto(userID, data)
	if err == nil && success {
		result.Success = success
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return result
}
