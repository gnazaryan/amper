package authorization

import (
	"amper/common/argument"
	"amper/common/structs"
	"amper/service/business"
)

// GetEntities is running a query to retrieve entities with the start and limit paging parameters
func GetEntities(userID *int64, start *int, limit *int) ([]structs.Entity, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "start": start, "limit": limit})
	if err != nil {
		return nil, err
	}
	//TODO perform authorization for get user action with userId
	entities, error := business.GetEntities(start, limit)
	return entities, error
}

// GetEntity is for getin an existing entity with give id
func GetEntity(userID *int64, EntityId *int64) (*structs.Entity, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "entityId": EntityId})
	if err != nil {
		return nil, err
	}
	return business.GetEntity(userID, EntityId)
}

// DeleteEntity is for deleting an existing entity with given id
func DeleteEntity(userID *int64, EntityId *int64) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "entityId": EntityId})
	if err != nil {
		return false, err
	}
	return business.DeleteEntity(userID, EntityId)
}

// CreateEntity is for creting an entity with the provided information
func CreateEntity(userID *int64, ApiName *string, Title *string, TitlePlural *string) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "apiName": ApiName, "title": Title, "titlePlural": TitlePlural})
	if err != nil {
		return false, err
	}
	return business.CreateEntity(userID, ApiName, Title, TitlePlural)
}

// EditEntity is for editing an entity with the id and provided information
func EditEntity(userID *int64, EntityId *int64, ApiName *string, Title *string, TitlePlural *string) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "entityId": EntityId, "apiName": ApiName, "title": Title, "titlePlural": TitlePlural})
	if err != nil {
		return false, err
	}
	return business.EditEntity(userID, EntityId, ApiName, Title, TitlePlural)
}

func GetFields(UserId *int64, ObjectId *int64) (*[]structs.Field, error) {
	err := argument.Validate(map[string]interface{}{"userID": UserId, "ObjectId": ObjectId})
	if err != nil {
		return nil, err
	}
	return business.GetFields(UserId, ObjectId)
}
