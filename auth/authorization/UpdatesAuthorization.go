package authorization

import (
	"amper/common/argument"
	"amper/service/business"
)

func FetchUpdates(userID *int64) (result map[string][]*interface{}, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID})
	if err != nil {
		return nil, err
	}
	return business.FetchUpdates(userID)
}

func PutUpdates(userID *int64, Category *string, Participants *[]int64, Value *interface{}) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "category": Category, "participants": Participants})
	if err != nil {
		return false, err
	}
	return business.PutUpdates(userID, Category, Participants, Value)
}
