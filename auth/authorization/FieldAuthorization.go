package authorization

import (
	"amper/common/argument"
	"amper/service/business"
)

func CreateField(userID *int64, entityId *int64, apiName *string, label *string, dataType *string, textLength *int64, objectReference *int64, required *bool, status *bool) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "entityId": entityId, "apiName": apiName, "label": label, "dataType": dataType, "required": required, "status": status})
	if err != nil {
		return false, err
	}
	return business.CreateField(userID, entityId, apiName, label, dataType, textLength, objectReference, required, status)
}

func DeleteField(userID *int64, entityId *int64, fieldIds *[]int64) (bool, error) {
	err := argument.Validate(map[string]interface{}{"userID": userID, "fieldIds": fieldIds, "entityId": entityId})
	if err != nil {
		return false, err
	}
	return business.DeleteField(userID, entityId, fieldIds)
}

func AddObjectTypeField(userID *int64, fieldId *int64, objectTypeId *int64) (bool, error) {
	err := argument.Validate(map[string]interface{}{"fieldId": fieldId, "objectTypeId": objectTypeId})
	if err != nil {
		return false, err
	}
	return business.AddObjectTypeField(userID, fieldId, objectTypeId)
}

func DeleteObjectTypeField(userID *int64, fieldId *int64, objectTypeId *int64) (bool, error) {
	err := argument.Validate(map[string]interface{}{"fieldId": fieldId, "objectTypeId": objectTypeId})
	if err != nil {
		return false, err
	}
	return business.DeleteObjectTypeField(userID, fieldId, objectTypeId)
}
