package controller

import (
	"amper/common/structs"
	"amper/common/util"
	"amper/common/util/arrays"
	"amper/service/business"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

// Dispatch comment
func Dispatch(w *http.ResponseWriter, r *http.Request) {
	result := "Welcome to Amper"
	if len(r.URL.Path) > 0 {

		success, userID, sessionId := itercept(w, r)
		if !success {
			unAuthenticatedResult, _ := json.Marshal(structs.Result{
				Success:       false,
				Error:         "We prefere authenticated users",
				Authenticated: -1,
			})
			(*w).Write([]byte(unAuthenticatedResult))
			return
		}

		controllerName := strings.Split(r.URL.Path, "/")[1]
		switch controllerName {
		case "/":
		case "amper":
			result = AmperController(&userID, sessionId, w, r)
		case "updates":
			result = UpdatesController(&userID, sessionId, w, r)
		case "users":
			result = UserController(&userID, sessionId, w, r)
		case "profile":
			result = ProfileController(&userID, sessionId, w, r)
		case "profiles":
			result = SecurityProfileController(&userID, w, r)
		case "email":
			result = EmailController(&userID, w, r)
		case "entities":
			result = EntityController(&userID, w, r)
		case "objectTypes":
			result = EntityTypeController(&userID, w, r)
		case "fields":
			result = FieldController(&userID, w, r)
		case "dashboards":
			result = DashboardController(&userID, w, r)
		case "widgets":
			result = WidgetController(&userID, w, r)
		case "records":
			result = RecordController(&userID, w, r)
		case "files":
			result = FileController(&userID, w, r)
		case "files-v1":
			result = FileV1Controller(&userID, w, r)
		case "settings":
			result = SettingsController(&userID, w, r)
		case "chat":
			result = ChatController(&userID, sessionId, w, r)
		default:
		}
	}
	if len(result) > 0 {
		(*w).Write([]byte(result))
	}
}

func getSessionId(r *http.Request) *string {
	sessionId := r.URL.Query().Get("sessionId")
	if util.EmptyString(&sessionId) {
		sessionId = r.Header.Get("sessionId")
		if util.EmptyString(&sessionId) {
			cookie, errC := r.Cookie("user_sessionId")
			if errC == nil {
				sessionId, _ = url.QueryUnescape(cookie.Value)
				if len(sessionId) > 0 && sessionId[0] == '"' {
					sessionId = sessionId[1:]
				}
				if len(sessionId) > 0 && sessionId[len(sessionId)-1] == '"' {
					sessionId = sessionId[:len(sessionId)-1]
				}
			}
			/*cookies := r.Cookies()
			for _, cookie := range cookies {
				name := cookie.Name
				value := cookie.Value
				if name == "user_sessionId" {
					sessionId = value
					break
				}
			}*/
		}
	}
	return &sessionId
}

var excludeUrls = []string{"/users/login", "/users/activate", "/amper/invalidateCache"}

func itercept(w *http.ResponseWriter, r *http.Request) (success bool, userId int64, sessionId *string) {
	if !arrays.Contains(excludeUrls, r.URL.Path) {
		sessionId := getSessionId(r)
		if !util.EmptyString(sessionId) {
			return business.ValidateSession(sessionId)
		}
	} else {
		return true, -1, nil
	}
	return false, -1, nil
}
