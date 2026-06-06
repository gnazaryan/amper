package authorization

import (
	"amper/common/argument"
	"amper/common/structs"
	"amper/service/business"
)

func GetInstances(userID *int64, Type *string) ([]structs.Amper, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID})
	if err != nil {
		return nil, err
	}
	//TODO perform authorization for get user action with userId
	instances, error := business.GetInstances(userID, Type)
	return instances, error
}

func RemoveInstance(userID *int64, amper structs.Amper) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "instanceId": amper.Id})
	if err != nil {
		return false, err
	}
	success, error := business.RemoveInstance(userID, amper)
	return success, error
}

func CreateInstance(userID *int64, sessionId *string, amper structs.Amper) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "sessionId": sessionId, "Identifier": amper.Identifier, "name": amper.Name, "type": amper.Type, "address": amper.Address, "port": amper.Port, "limit": amper.Limit, "directory": amper.Directory})
	if err != nil {
		return false, err
	}
	success, error := business.CreateInstance(userID, sessionId, amper)
	return success, error
}

func EditInstance(userID *int64, sessionId *string, amper structs.Amper) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "sessionId": sessionId, "name": amper.Name, "address": amper.Address, "port": amper.Port, "limit": amper.Limit, "directory": amper.Directory})
	if err != nil {
		return false, err
	}
	success, error := business.EditInstance(userID, sessionId, amper)
	return success, error
}

func FetchInstanceInfo(userID *int64, sessionId *string, amper structs.Amper) (*structs.Amper, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "type": amper.Type, "address": amper.Address, "port": amper.Port})
	if err != nil {
		return nil, err
	}
	info, error := business.FetchInstanceInfo(userID, sessionId, amper)
	return info, error
}

func FetchStatus(userID *int64, amper structs.Amper) (*structs.Amper, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "directory": amper.Directory})
	if err != nil {
		return nil, err
	}
	info, err := business.FetchStatus(userID, amper)
	return info, err
}

func InvalidateCache(userID *int64, Name *string, userIdDelete *int64, chatChannelDelete *int64) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "name": Name})
	if err != nil {
		return false, err
	}
	success, err := business.InvalidateCache(userID, Name, userIdDelete, chatChannelDelete)
	return success, err
}
