package controller

import (
	"amper/auth/authorization"
	"amper/common/structs"
	"amper/common/util"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// UserController is responsible for dispatching requests related to
// user managment functionalities
func UserController(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	var resultStruct interface{}
	if len(pathSplit) > 2 {
		switch pathSplit[2] {
		case "login":
			resultStruct = login(userID, w, r)
		case "fetch":
			resultStruct = getUsers(userID, w, r)
		case "create":
			resultStruct = create(userID, w, r)
		case "download":
		case "edit":
			resultStruct = edit(userID, sessionId, w, r)
		case "getActivatingUser":
		case "activate":
			resultStruct = activate(userID, w, r)
		case "isValidUserName":
			resultStruct = isValidUserName(userID, w, r)
		case "remove":
			resultStruct = remove(userID, w, r)
		case "getUserRelationships":
			resultStruct = getUserRelationships(userID, w, r)
		case "createUserRelationship":
			resultStruct = createUserRelationship(userID, w, r)
		case "deleteUserRelationship":
			resultStruct = deleteUserRelationship(userID, w, r)
		default:
		}
		marshaled, _ := json.Marshal(resultStruct)
		result = string(marshaled)
	}
	return
}

func remove(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ValidationResult) {
	var parameters struct {
		UserID *int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	valid, err := authorization.RemoveSoft(userID, parameters.UserID)
	if err == nil {
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	result.Valid = valid
	return
}

func isValidUserName(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ValidationResult) {
	var parameters struct {
		Username *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	valid, err := authorization.IsValidUserName(userID, parameters.Username)
	if err == nil {
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	result.Valid = valid
	return
}

func activate(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		ActivationCode  *string
		Password        *string
		ConfirmPassword *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, err := authorization.Activate(parameters.ActivationCode, parameters.Password)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return
}

func edit(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var user structs.User
	json.NewDecoder(r.Body).Decode(&user)
	success, err := authorization.EditUser(userID, sessionId, user)
	result.Success = success
	if err != nil {
		result.Error = err.Error()
	}
	return
}

// create is responsible for creating a new user with the specified parametes
func create(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var user structs.User
	json.NewDecoder(r.Body).Decode(&user)
	success, err := authorization.CreateUser(userID, user)
	if err == nil {
		result.Success = success
	} else {
		result.Error = err.Error()
	}
	return
}

// getUsers is resposible for retrieving users with the provide start and limit parameters and the
// user requesting the information should have permissions
func getUsers(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.UserAndProfilesReslt) {
	var parameters struct {
		Start         *int
		Limit         *int
		Search        *[]string
		SorfField     *string
		SortDirection *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	userAndProfiles, totalCount, err := authorization.GetUsers(userID, parameters.Start, parameters.Limit, parameters.Search, parameters.SorfField, parameters.SortDirection)
	if err == nil {
		result.Data = userAndProfiles
		result.TotalCount = totalCount
	} else {
		result.Error = err.Error()
	}
	return
}

// userLogin is responsible for authorizing user with the provided username and password
// and returning session id if the authentication was successfull
func login(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.UserLoginReslt) {
	var parameters struct {
		Username string
		Password string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	user, err := authorization.UserLogin(&parameters.Username, &parameters.Password)
	if user != nil && !util.EmptyString(&user.SessionID) && err == nil {
		settings := GetSettings(userID)
		setCookies(userID, w, user, settings)
		result.SessionID = user.SessionID
		result.User = user
		result.Settings = GetSettings(userID)
		result.Success = true
	} else {
		result.Error = err.Error()
	}
	return
}

func setCookies(userID *int64, w *http.ResponseWriter, user *structs.User, settings *structs.Settings) {
	amper := AmperInfo(userID)
	setCookie(userID, *amper.Data, w, "sessionId", user.SessionID)
}

func setCookie(userID *int64, instance structs.Amper, w *http.ResponseWriter, Name string, Value string) {
	http.SetCookie(*w, &http.Cookie{
		Name:     fmt.Sprintf("%s_%d", Name, instance.Id),
		Value:    Value,
		MaxAge:   3 * 60 * 60,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})
}

func GetSettings(userID *int64) *structs.Settings {
	adobeLicenseKey := authorization.GetSetting(userID, util.PointerString("amper.adobeLicenseKey"), util.PointerString(""))
	settings := structs.Settings{
		AdobeLicenseKey: &adobeLicenseKey,
	}
	return &settings
}

// getUserRelationships is resposible for retrieving user relationships with the provide employee id parameter
func getUserRelationships(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.UserRelationshipExtendedResult) {
	var parameters struct {
		EmployeeId int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	userRelationships, err := authorization.GetUserRelationships(*userID, parameters.EmployeeId)
	if err == nil {
		result.Data = userRelationships
	} else {
		result.Error = err.Error()
	}
	return
}

func createUserRelationship(userId *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		EmployeeId int64
		ManagerId  int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	resultCUR, err := authorization.CreateUserRelationship(*userId, parameters.EmployeeId, parameters.ManagerId)
	if err == nil {
		result.Success = resultCUR
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return
}

func deleteUserRelationship(userId *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		ManagerId  int64
		EmployeeId int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	resultCUR, err := authorization.DeleteUserRelationship(*userId, parameters.ManagerId, parameters.EmployeeId)
	if err == nil {
		result.Success = resultCUR
	} else {
		result.Error = err.Error()
		result.Success = false
	}
	return
}
