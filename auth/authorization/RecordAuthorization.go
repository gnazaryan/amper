package authorization

import (
	"amper/common/argument"
	"amper/common/structs"
	"amper/service/business"
	"errors"
)

func FetchRecords(userId *int64, apiName *string, objectId *int64, start *int64, limit *int64, search *string, metadata *bool) (result *[]map[string]interface{}, metadadata *structs.ObjectMetadata, totalCount int64, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "start": start, "limit": limit})
	if err != nil {
		return nil, nil, 0, err
	}
	if apiName == nil && objectId == nil {
		return nil, nil, 0, errors.New("either apiName or objectId must be supplied with values")
	}
	result, metadadata, totalCount, err = business.FetchRecords(userId, apiName, objectId, start, limit, search, metadata)
	return result, metadadata, totalCount, err
}

func AddRecord(userId *int64, apiName *string, payload *string) (result *structs.Record, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "apiName": apiName, "payload": payload})
	if err != nil {
		return nil, err
	}
	result, err = business.AddRecord(userId, apiName, payload)
	return result, err
}

func RemoveRecord(userId *int64, apiName *string, identifier *string, id *int64) (result *structs.Record, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "identifier": identifier})
	if err != nil {
		return nil, err
	}
	result, err = business.RemoveRecord(userId, apiName, identifier, id)
	return result, err
}

func UpdateRecord(userId *int64, apiName *string, payload *string) (result *structs.Record, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "payload": payload})
	if err != nil {
		return nil, err
	}
	result, err = business.UpdateRecord(userId, apiName, payload)
	return result, err
}

func AddRecords(userId *int64, apiName *string, payload *string) (resultSuccess *[]map[string]interface{}, resultError *[]map[string]interface{}, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "apiName": apiName, "payload": payload})
	if err != nil {
		return nil, nil, err
	}
	resultSuccess, resultError, err = business.AddRecords(userId, apiName, payload)
	return resultSuccess, resultError, err
}

func RemoveRecords(userId *int64, apiName *string, ids *[]int64, identifiers *[]string) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "identifiers": identifiers})
	if err != nil {
		return false, err
	}
	result, err = business.RemoveRecords(userId, apiName, ids, identifiers)
	return result, err
}

func UpdateRecords(userId *int64, payloads *string) (resultSuccess *[]map[string]interface{}, resultError *[]map[string]interface{}, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "payloads": payloads})
	if err != nil {
		return nil, nil, err
	}
	resultSucc, resultErr, err := business.UpdateRecords(userId, payloads)
	resultSuccess = &resultSucc
	resultError = &resultErr
	return
}
