package authorization

import (
	"amper/common/argument"
	"amper/common/structs"
	"amper/service/business"
)

// GetObjectTypes is running a query to retrieve entitiy types
func GetObjectTypes(UserID *int64, ObjectId *int64) (result *[]structs.ObjectType, err error) {
	err = argument.Validate(map[string]interface{}{"userID": UserID, "objectId": ObjectId})
	if err != nil {
		return nil, err
	}
	result, err = business.GetObjectTypes(UserID, ObjectId)
	return
}

func CreateObjectType(UserID *int64, ObjectId *int64, ExtendsTo *int64, ApiName *string, Label *string) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": UserID, "objectId": ObjectId, "ExtendsTo": ExtendsTo, "ApiName": ApiName, "Label": Label})
	if err != nil {
		return false, err
	}
	result, err = business.CreateObjectType(UserID, ObjectId, ExtendsTo, ApiName, Label)
	return
}

func DeleteObjectType(UserID *int64, ObjectTypeId *int64) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": UserID, "ObjectTypeId": ObjectTypeId})
	if err != nil {
		return false, err
	}
	result, err = business.DeleteObjectType(UserID, ObjectTypeId)
	return
}

// GetObjectTypes is running a query to retrieve entitiy types
func GetObjectTypeFields(UserID *int64, ObjectId *int64, ObjectTypeId *int64) (result *structs.ObjectType, err error) {
	err = argument.Validate(map[string]interface{}{"userID": UserID, "objectId": ObjectId, "objectTypeId": ObjectTypeId})
	if err != nil {
		return nil, err
	}
	result, err = business.GetObjectTypeFields(UserID, ObjectId, ObjectTypeId)
	return
}
