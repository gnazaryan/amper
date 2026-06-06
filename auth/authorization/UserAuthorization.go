package authorization

import (
	"amper/common/argument"
	"amper/common/structs"
	"amper/service/business"
)

// UserLogin is responsible for authenticating and creating a session for the
// username and password
func UserLogin(username *string, password *string) (*structs.User, error) {
	err := argument.Validate(map[string]interface{}{"username": username, "password": password})
	if err != nil {
		return nil, err
	}
	//No authorization required for user login
	user, error := business.UserLogin(username, password)
	return user, error
}

// GetUsers is running a query to retrieve users with the start and limit paging parameters
func GetUsers(userID *int64, start *int, limit *int, search *[]string, sortField *string, sortDirection *string) ([]structs.UserAndProfile, int, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "start": start, "limit": limit})
	if err != nil {
		return nil, 0, err
	}
	//TODO perform authorization for get user action with userId
	users, totalCount, errGU := business.GetUsers(start, limit, search, sortField, sortDirection)
	return users, totalCount, errGU
}

// CreateUser is for creating/registering a new user
func CreateUser(userID *int64, user structs.User) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "email": user.Email, "firstName": user.FirstName,
		"lastName": user.LastName, "profile": user.Profile, "username": user.Username, "amperId": user.AmperId})
	if err != nil {
		return false, err
	}
	return business.CreateUser(userID, user)
}

// EditUser is for editing/pdating an existing user
func EditUser(userID *int64, sessionId *string, user structs.User) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "sessionId": sessionId, "id": user.ID, "email": user.Email, "firstName": user.FirstName,
		"lastName": user.LastName, "profile": user.Profile, "username": user.Username})
	if err != nil {
		return false, err
	}
	return business.EditUser(userID, sessionId, user)
}

// Activate is for activating a user with the provided activation code
func Activate(activationCode *string, password *string) (bool, error) {
	err := argument.Validate(map[string]interface{}{"activationCode": activationCode, "password": password})
	if err != nil {
		return false, err
	}
	return business.Activate(activationCode, password)
}

// IsValidUserName is for checking if the username is avalable
func IsValidUserName(userID *int64, username *string) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "username": username})
	if err != nil {
		return false, err
	}
	return business.IsValidUserName(userID, username)
}

// Remove is for deleting a user with the given user id
func Remove(userID *int64, userIDToRemove *int64) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "userId": userIDToRemove})
	if err != nil {
		return false, err
	}
	return business.Remove(userID, userIDToRemove)
}

// Remove is for deleting a user with the given user id
func RemoveSoft(userID *int64, userIDToRemove *int64) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "userId": userIDToRemove})
	if err != nil {
		return false, err
	}
	return business.RemoveSoft(userID, userIDToRemove)
}
