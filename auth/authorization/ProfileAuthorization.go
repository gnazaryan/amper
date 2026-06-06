package authorization

import (
	"amper/common/argument"
	"amper/common/structs"
	"amper/service/business"
	"io"
)

// GetUsers is running a query to retrieve users with the start and limit paging parameters
func UpdateCover(userID *int64, data []byte) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "data": data})
	if err != nil {
		return false, err
	}
	//TODO perform authorization for get user action with userId
	success, error := business.UpdateCover(userID, data)
	return success, error
}

func ViewCover(userID *int64) (result *io.ReadCloser, metadata *structs.FileMetadata, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID})
	if err != nil {
		return nil, nil, err
	}
	return business.ViewCover(userID)
}

// GetUsers is running a query to retrieve users with the start and limit paging parameters
func UpdatePhoto(userID *int64, data []byte) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "data": data})
	if err != nil {
		return false, err
	}
	//TODO perform authorization for get user action with userId
	success, error := business.UpdatePhoto(userID, data)
	return success, error
}

func ViewPhoto(userID *int64) (result *io.ReadCloser, metadata *structs.FileMetadata, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID})
	if err != nil {
		return nil, nil, err
	}
	return business.ViewPhoto(userID)
}

func AdjustPhoto(userID *int64, sessionId *string, PositionX *int, PositionY *int, Width *int, Height *int) (success bool, result *string, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "PositionX": PositionX, "PositionY": PositionY, "Width": *Width, "Height": *Height})
	if err != nil {
		return false, nil, err
	}
	return business.AdjustPhoto(userID, sessionId, PositionX, PositionY, Width, Height)
}

func GetProfileState(userID *int64) (data map[string]interface{}, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID})
	if err != nil {
		return nil, err
	}
	return business.GetProfileState(userID)
}

func SaveConfiguration(userId *int64, Name *string, Value *map[string]interface{}) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userId, "Name": Name})
	if err != nil {
		return false, err
	}
	return business.SaveConfiguration(userId, Name, Value)
}

func SaveDetail(userId *int64, Name *string, Value *string) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userId, "Name": Name})
	if err != nil {
		return false, err
	}
	return business.SaveDetail(userId, Name, Value)
}

func AddRelationship(userID *int64, employeeId int64, Type *string, Value int64) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "employeeId": employeeId, "Type": Type, "Value": Value})
	if err != nil {
		return false, err
	}
	return business.AddRelationship(userID, employeeId, Type, Value)
}

func RemoveRelationship(userID *int64, employeeId int64, Type *string, Value int64) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "employeeId": employeeId, "Type": Type, "Value": Value})
	if err != nil {
		return false, err
	}
	return business.RemoveRelationship(userID, employeeId, Type, Value)
}
