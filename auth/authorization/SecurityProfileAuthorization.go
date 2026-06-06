package authorization

import (
	"amper/common/argument"
	"amper/common/structs"
	"amper/service/business"
)

// GetUsers is running a query to retrieve users with the start and limit paging parameters
func FetchProfiles(userID *int64, start *int, limit *int, search *[]string, sortField *string, sortDirection *string) ([]structs.Profile, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "start": start, "limit": limit})
	if err != nil {
		return nil, err
	}
	//TODO perform authorization for get user action with userId
	profiles, error := business.FetchProfiles(start, limit, search, sortField, sortDirection)
	return profiles, error
}
