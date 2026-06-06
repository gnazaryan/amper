package structs

import (
	"amper/common/util"
	"amper/common/util/jsons"
	"encoding/json"
	"fmt"
	"time"
)

// Session struct representing user authentication
type Session struct {
	SessionID string `json:"sessionId"`
	DateTime  time.Time
	UserID    int64
}

// UserAndProfile struct representing user and profile
type UserAndProfile struct {
	User
	ProfileID   int    `json:"profileId"`
	ProfileName string `json:"profileName"`
}

// UserAndProfilesReslt struct representing users
type UserAndProfilesReslt struct {
	Result
	Data       []UserAndProfile `json:"data"`
	TotalCount int              `json:"totalCount"`
}

// UserAndProfilesReslt struct representing users
type UserReslt struct {
	Result
	Data       []User `json:"data"`
	TotalCount int    `json:"totalCount"`
}

// User struct representing user
type User struct {
	Session
	ID             *int64                    `json:"id"`
	FirstName      *string                   `json:"firstName"`
	LastName       *string                   `json:"lastName"`
	MiddleName     *string                   `json:"middleName"`
	Username       *string                   `json:"username"`
	Password       *string                   `json:"password"`
	Photo          *string                   `json:"photo"`
	Profile        *int64                    `json:"profileId"`
	Email          *string                   `json:"email"`
	Active         *int                      `json:"active"`
	ActivationCode *string                   `json:"activationCode"`
	AmperId        *int64                    `json:"amperId"`
	State          *int                      `json:"state"`
	Config         *string                   `json:"config"`
	Emails         *[]map[string]interface{} `json:"emails"`
}

func (u *User) Initialize(sensitive bool) error {
	var config map[string]interface{}
	if u.Config != nil {
		var errJ error
		config, errJ = jsons.GetJsonObject(u.Config)
		if errJ != nil {
			util.Loggify(errJ)
			return fmt.Errorf("not able to parse json config for the user")
		}
	} else {
		config = make(map[string]interface{})
	}
	settings, okS := config["settings"].(map[string]interface{})
	if okS {
		email, okE := settings["email"].([]interface{})
		if okE {
			emails := make([]map[string]interface{}, 0)
			u.Emails = &emails
			for _, item := range email {
				emailItem := make(map[string]interface{})
				email, okEm := item.(map[string]interface{})["email"].(string)
				if okEm {
					emailItem["email"] = email
				}
				label, okL := item.(map[string]interface{})["label"].(string)
				if okL {
					emailItem["label"] = label
				}
				mailboxes, okM := item.(map[string]interface{})["mailboxes"].([]interface{})
				if okM {
					emailItem["mailboxes"] = mailboxes
				}
				if sensitive {
					password, okP := item.(map[string]interface{})["password"].(string)
					if okP {
						emailItem["password"] = password
					}
				}
				emails = append(emails, emailItem)
			}
		}
	}
	return nil
}

func (u *User) AddConfig(key string, value map[string]interface{}) error {
	var config map[string]interface{}
	if u.Config != nil {
		var errJ error
		config, errJ = jsons.GetJsonObject(u.Config)
		if errJ != nil {
			util.Loggify(errJ)
			return fmt.Errorf("not able to parse json config for the user")
		}
	} else {
		config = make(map[string]interface{})
	}
	config[key] = value
	resultM, errJM := json.Marshal(config)
	if errJM != nil {
		util.Loggify(errJM)
		return fmt.Errorf("not able to marshal map to json string")
	}
	u.Config = util.PointerString(string(resultM))
	return nil
}

func (u *User) GetConfig() (result map[string]interface{}, err error) {
	if u.Config != nil {
		var errJ error
		result, errJ = jsons.GetJsonObject(u.Config)
		if errJ != nil {
			util.Loggify(errJ)
			return nil, fmt.Errorf("not able to parse json config for the user")
		}
	} else {
		result = make(map[string]interface{})
	}
	return result, nil
}

// Entity struct representing entity
type Entity struct {
	Session
	ID          *int64  `json:"id"`
	ApiName     *string `json:"apiName"`
	Title       *string `json:"title"`
	TitlePlural *string `json:"titlePlural"`
}

// Entities struct representing entities
type Entities struct {
	Result
	Data []Entity `json:"data"`
}

// EntityResponse struct representing entities
type EntityResponse struct {
	Result
	Entity *Entity `json:"entity"`
}

// UserLoginReslt struct representing login result
type UserLoginReslt struct {
	Result
	SessionID string    `json:"sessionId"`
	User      *User     `json:"user"`
	Settings  *Settings `json:"settings"`
}

// UserNotification struct representing the notification tempalte values for user registration
type UserNotification struct {
	UserFirstName *string
	ButtonLabel   *string
	ButtonHref    *string
}

// The user profile data model
type Profile struct {
	ID   *int64  `json:"id"`
	Name *string `json:"name"`
}

// The http result for the user profiles
type ProfileResult struct {
	Result
	Data []Profile `json:"data"`
}

type UserRelationship struct {
	ID         *int64 `json:"id"`
	EmployeeId *int64 `json:"employeeId"`
	ManagerId  *int64 `json:"managerId"`
}

type UserRelationshipExtended struct {
	ID                *int64  `json:"id"`
	EmployeeId        *int64  `json:"employeeId"`
	EmployeeFirstName *string `json:"employeeFirstName"`
	EmployeeLastName  *string `json:"employeeLastName"`
	EmployeePhoto     *string `json:"employeePhoto"`
	ManagerId         *int64  `json:"managerId"`
	ManagerFirstName  *string `json:"managerFirstName"`
	ManagerLastName   *string `json:"managerLastName"`
	ManagerPhoto      *string `json:"managerPhoto"`
}

type UserRelationshipExtendedResult struct {
	Result
	Data []UserRelationshipExtended `json:"data"`
}

type UserDetail struct {
	ID               *int64  `json:"id"`
	UserId           *int64  `json:"userId"`
	Info             *string `json:"info"`
	AboutMe          *string `json:"aboutMe"`
	Responsibilities *string `json:"responsibilities"`
	Skills           *string `json:"skills"`
}

func (u *UserDetail) GetInfo() (result map[string]interface{}, err error) {
	if u.Info != nil {
		var errJ error
		result, errJ = jsons.GetJsonObject(u.Info)
		if errJ != nil {
			util.Loggify(errJ)
			return nil, fmt.Errorf("not able to parse json config for the user")
		}
	} else {
		result = make(map[string]interface{})
	}
	return result, nil
}
