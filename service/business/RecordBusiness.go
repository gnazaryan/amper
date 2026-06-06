package business

import (
	"amper/cache/business"
	"amper/common/constants"
	"amper/common/structs"
	"amper/common/util"
	"amper/common/util/ampstrings"
	"amper/common/util/datetime"
	"amper/common/util/jsons"
	"amper/common/util/maps"
	"amper/data/database"
	"container/list"
	"encoding/base64"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func FetchRecords(userId *int64, apiName *string, objectId *int64, start *int64, limit *int64, search *string, metadata *bool) (result *[]map[string]interface{}, metadadata *structs.ObjectMetadata, totalCount int64, err error) {
	var searchParams map[string]interface{}
	if search != nil && len(*search) > 0 {
		var errJ error
		searchParams, errJ = jsons.GetJsonObject(search)
		if errJ != nil {
			log.Print(errJ.Error(), errJ)
			err = fmt.Errorf("unable to fetch records, the supplied search data is in a wrong format")
			return
		}
	}
	var searchParameter = structs.Search{}
	if !searchParameter.Parse(search) {
		err = fmt.Errorf(strings.Join(*searchParameter.Errors, ", "))
		return
	}

	var object *structs.Entity
	if apiName != nil && len(*apiName) > 0 {
		var errO error
		object, errO = GetEntityByApiName(userId, apiName)
		if errO != nil {
			log.Print(errO.Error(), errO)
			err = fmt.Errorf("unable to fetch records, the supplied api name %s is not recongnized", *apiName)
			return
		}
	} else if objectId != nil {
		var errO error
		object, errO = GetEntity(userId, objectId)
		if errO != nil {
			log.Print(errO.Error(), errO)
			err = fmt.Errorf("unable to fetch records, the supplied object id %d is not recongnized", *objectId)
			return
		}
	} else {
		err = fmt.Errorf("unable to fetch records, either object if or api name must be supplied with values")
		return
	}

	fields, errF := GetFields(userId, object.ID)
	if errF != nil {
		log.Print(errF.Error(), errF)
		err = fmt.Errorf("unable to fetch records, not able to retrieve field metadata information for object %s", *apiName)
		return
	}
	foreginKey := make(map[string]string)
	for _, field := range *fields {
		if constants.GetDataTypes.REFERENCE.Name == *field.Type {
			objectReferenceId := field.ObjectReference
			object, errO := GetEntity(userId, objectReferenceId)
			if errO == nil && object != nil && object.ID != nil {
				foreginKey[*field.ApiName] = *object.ApiName
			} else {
				//Fail silently, no need to bubble up
				util.Loggify(errO)
			}
		}
	}

	if metadata != nil && *metadata {
		objectTypes, errOT := GetObjectTypes(userId, object.ID)
		if errOT != nil {
			log.Print(errOT.Error(), errOT)
			err = fmt.Errorf("unable to fetch records, not able to retrieve object type metadata information for object %s", *apiName)
		}
		metadadata = &structs.ObjectMetadata{
			Object:      object,
			ObjectTypes: objectTypes,
			Fields:      fields,
		}
	}
	result, totalCount, errDB := database.FetchRecords(userId, object.ApiName, start, limit, searchParams, &searchParameter, &foreginKey)
	if errDB != nil {
		log.Print(errDB.Error(), errDB)
		err = fmt.Errorf("unable to fetch records, please try again later or contact the support")
	}
	return result, metadadata, totalCount, err
}

func AddRecord(userId *int64, apiName *string, payload *string) (result *structs.Record, err error) {
	objectMetadata, err := GetObjectMetadata(userId, apiName, nil, nil, nil)
	if err != nil {
		return
	}
	record, errJ := jsons.GetJsonObject(payload)
	if errJ != nil {
		err = fmt.Errorf("unable to add a record, not able to parse the supplied payload, please fix and try again or contact the support")
		return
	}
	objectTypeKey := constants.GetBaseObjectType().KEY
	if record[constants.GetDefaultFields().OBJECTTYPE.ApiName] != nil {
		if reflect.TypeOf(record[constants.GetDefaultFields().OBJECTTYPE.ApiName]).String() != "string" {
			err = fmt.Errorf("unable to add a record, the supplied object type info is wrong, please fix and try again or contact the support")
			return
		} else {
			objectTypeKey = record[constants.GetDefaultFields().OBJECTTYPE.ApiName].(string)
		}
	}
	var objectType *structs.ObjectType
	for _, value := range *objectMetadata.ObjectTypes {
		if *value.ApiName == objectTypeKey {
			objectType = &value
			break
		}
	}
	if objectType == nil {
		err = fmt.Errorf("unable to add a record, not able to locate the object type info, please fix and try again or contact the support")
		return
	}
	fieldsMap := fieldsToMap(objectMetadata.Fields)
	var dbRecord = make(map[string]string)
	var errorValues = make([]string, 0)
	for _, value := range objectType.ObjectTypeFields {
		field := fieldsMap[*value.FieldId]
		if field.Status == nil || *field.Status != 1 || *value.ApiName == constants.GetDefaultFields().ID.ApiName || *value.ApiName == constants.GetDefaultFields().IDENTIFIER.ApiName || *value.ApiName == constants.GetDefaultFields().OBJECTTYPE.ApiName {
			continue
		}
		if record[*field.ApiName] == nil {
			if *field.Required > 0 {
				errorValues = append(errorValues, fmt.Sprintf("The %s is a required field, it must be supplied with a value of type %s", *field.ApiName, *field.Type))
			}
		} else {
			dirtyValue := record[*field.ApiName]
			convertedValue, errC := ConvertFieldValue(&dirtyValue, &field)
			if errC != nil {
				if *field.Required > 0 {
					errorValues = append(errorValues, errC.Error())
				}
			} else {
				dbRecord[*field.ApiName] = convertedValue
			}
		}
	}
	if len(errorValues) > 0 {
		return nil, fmt.Errorf(strings.Join(errorValues, ". "))
	}
	dbRecord[constants.GetDefaultFields().OBJECTTYPE.ApiName] = *objectType.ApiName
	dbRecord[constants.GetDefaultFields().IDENTIFIER.ApiName] = strconv.FormatInt(*business.AmperId(), 10) + "|" + strconv.FormatInt(*objectMetadata.Object.ID, 10) + "|" + strconv.FormatInt(*objectType.ID, 10)
	recordDb, errDb := database.AddRecord(userId, apiName, &dbRecord)
	if errDb != nil {
		log.Print(errDb.Error(), errDb)
		err = fmt.Errorf("unable to add a record, please try again later or contact the support")
		return
	}
	var dbRecordUpdate = make(map[string]string)
	dbRecordUpdate[constants.GetDefaultFields().ID.ApiName] = strconv.FormatInt(recordDb.ID, 10)
	dbRecordUpdate[constants.GetDefaultFields().IDENTIFIER.ApiName] = base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(*business.AmperId(), 10) + "|" + strconv.FormatInt(*objectMetadata.Object.ID, 10) + "|" + strconv.FormatInt(*objectType.ID, 10) + "|" + strconv.FormatInt(recordDb.ID, 10)))
	_, errDbU := database.UpdateRecord(userId, apiName, &dbRecordUpdate)
	if errDbU != nil {
		log.Print(errDbU.Error(), errDbU)
		err = fmt.Errorf("unable to add a record, please try again later or contact the support")
		_, errDbUR := database.RemoveRecord(userId, apiName, &recordDb.ID, nil)
		if errDbUR != nil {
			log.Print(errDbUR.Error(), errDbUR)
			err = fmt.Errorf("unable to add a record, wrong data created, please reach the support")
		}
	} else {
		recordDb.Record[constants.GetDefaultFields().IDENTIFIER.ApiName] = dbRecordUpdate[constants.GetDefaultFields().IDENTIFIER.ApiName]
		result = recordDb
	}
	return result, err
}

func ConvertFieldValue(intput *interface{}, field *structs.Field) (result string, err error) {
	switch *field.Type {
	case constants.GetDataTypes.TEXT.Name:
		if reflect.TypeOf(*intput).String() == "string" {
			length := int64(len((*intput).(string)))
			if length > *field.TextLength {
				err = fmt.Errorf("the supplied value for field %s has a longer content %d  then it was deseigned %d", *field.ApiName, length, *field.TextLength)
			} else if *field.Required == int8(1) && length < 1 {
				err = fmt.Errorf("the supplied value for field %s has no content, it is a required field and must contain value", *field.ApiName)
			} else {
				result = (*intput).(string)
			}
		} else {
			err = fmt.Errorf("the supplied value for field %s has a wrong data format", *field.ApiName)
		}
	case constants.GetDataTypes.BOOLEAN.Name:
		if reflect.TypeOf(*intput).String() == "bool" {
			if (*intput).(bool) {
				result = "1"
			} else {
				result = "0"
			}
		} else if reflect.TypeOf(*intput).String() == "string" {
			stringInput := (*intput).(string)
			if strings.EqualFold(stringInput, "true") || strings.EqualFold(stringInput, "1") || strings.EqualFold(stringInput, "yes") || strings.EqualFold(stringInput, "active") {
				result = "1"
			} else {
				result = "0"
			}
		} else if value, errV := util.I2Num(*intput); errV == nil && (value == 1 || value == 0) {
			result = strconv.FormatInt(value, 10)
		} else {
			err = fmt.Errorf("the supplied value for field %s has a wrong data format", *field.ApiName)
		}
	case constants.GetDataTypes.DATE.Name:
		if reflect.TypeOf(*intput).String() == "string" {
			stringInput := (*intput).(string)
			date, errD := datetime.ParseDate(&stringInput)
			if errD != nil {
				err = fmt.Errorf("the supplied value for field %s has a wrong date format, use %s format instead", *field.ApiName, datetime.DATE_FORMATS[0])
			} else {
				result = datetime.FormatDate(date)
			}
		} else {
			err = fmt.Errorf("the supplied value for field %s has a wrong date format, use %s format instead", *field.ApiName, datetime.DATE_FORMATS[0])
		}
	case constants.GetDataTypes.DATETIME.Name:
		if reflect.TypeOf(*intput).String() == "string" {
			stringInput := (*intput).(string)
			date, errD := datetime.ParseDateTime(&stringInput)
			if errD != nil {
				err = fmt.Errorf("the supplied value for field %s has a wrong date time format, use %s format instead", *field.ApiName, datetime.DATE_TIME_FORMATS[0])
			} else {
				result = datetime.FormatDateTime(date)
			}
		} else {
			err = fmt.Errorf("the supplied value for field %s has a wrong date time format, use %s format instead", *field.ApiName, datetime.DATE_TIME_FORMATS[0])
		}
	case constants.GetDataTypes.NUMBER.Name:
		if reflect.TypeOf(*intput).String() == "string" {
			stringInput := (*intput).(string)
			int64Input, errI := strconv.ParseInt(stringInput, 10, 64)
			if errI != nil {
				err = fmt.Errorf("the supplied value for field %s has a wrong number format", *field.ApiName)
			} else if int64Input > *field.TextLength {
				err = fmt.Errorf("the supplied value for field %s has larger value then it is expected %d", *field.ApiName, *field.TextLength)
			} else {
				result = strconv.FormatInt(int64Input, 10)
			}
		} else if strings.Contains(reflect.TypeOf(*intput).String(), "int") {
			int64Input, errI := util.I2Num(*intput)
			if errI != nil {
				err = fmt.Errorf("the supplied value for field %s has a wrong number format", *field.ApiName)
			} else if int64Input > *field.TextLength {
				err = fmt.Errorf("the supplied value for field %s has larger value then it is expected %d", *field.ApiName, *field.TextLength)
			} else {
				result = strconv.FormatInt(int64Input, 10)
			}
		} else {
			err = fmt.Errorf("the supplied value for field %s has a wrong number format, use number instead", *field.ApiName)
		}
	case constants.GetDataTypes.REFERENCE.Name:
		if reflect.TypeOf(*intput).String() == "string" {
			stringInput := (*intput).(string)
			decoded, errB := base64.StdEncoding.DecodeString(stringInput)
			if errB != nil {
				err = fmt.Errorf("the supplied value for field %s has a wrong reference identifier format", *field.ApiName)
			} else {
				idParts := strings.Split(string(decoded), "|")
				if len(idParts) == 4 {
					_, errI := strconv.ParseInt(idParts[3], 10, 64)
					if errI != nil {
						err = fmt.Errorf("the supplied value for field %s has a wrong reference  identifier format", *field.ApiName)
					} else {
						result = stringInput
					}
				} else {
					err = fmt.Errorf("the supplied value for field %s has a wrong reference identifier format", *field.ApiName)
				}
			}
		} else {
			err = fmt.Errorf("the supplied value for field %s has a wrong reference format, use valid reference instead", *field.ApiName)
		}
	}
	return result, err
}

func fieldsToMap(fields *[]structs.Field) (result map[int64]structs.Field) {
	result = make(map[int64]structs.Field)
	for _, value := range *fields {
		result[*value.ID] = value
	}
	return result
}

func objectTypesToMap(objectTypes *[]structs.ObjectType) (result map[string]*structs.ObjectType) {
	result = make(map[string]*structs.ObjectType)
	for i := 0; i < len(*objectTypes); i++ {
		ojectType := (*objectTypes)[i]
		result[*ojectType.ApiName] = &ojectType
	}
	return result
}

func GetObjectMetadata(userId *int64, apiName *string, objectProvided *structs.Entity, fieldsProvided *[]structs.Field, objectTypeProvided *[]structs.ObjectType) (result structs.ObjectMetadata, err error) {
	if objectProvided == nil {
		object, errO := GetEntityByApiName(userId, apiName)
		if errO != nil {
			log.Print(errO.Error(), errO)
			err = fmt.Errorf("unable to add a record as not able to find the object %s, please try again later or contact the support", *apiName)
			return
		} else if object == nil || object.ID == nil {
			err = fmt.Errorf("unable to add a record as not able to find the object %s, please try again later or contact the support", *apiName)
			return
		}
		result.Object = object
	} else {
		result.Object = objectProvided
	}
	if fieldsProvided == nil {
		fields, errF := GetFields(userId, result.Object.ID)
		if errF != nil {
			log.Print(errF.Error(), errF)
			err = fmt.Errorf("unable to add a record, not able to find the fields for object %s, please try again later or contact the support", *apiName)
			return
		} else if fields == nil || len(*fields) < 1 {
			err = fmt.Errorf("unable to add a record, not able to find the fields for object %s, please try again later or contact the support", *apiName)
			return
		}
		result.Fields = fields
	} else {
		result.Fields = fieldsProvided
	}

	if objectTypeProvided == nil {
		objectTypes, errOT := GetObjectTypes(userId, result.Object.ID)
		if errOT != nil {
			log.Print(errOT.Error(), errOT)
			err = fmt.Errorf("unable to add a record, not able to find the object types for object %s, please try again later or contact the support", *apiName)
			return
		} else if objectTypes == nil || len(*objectTypes) < 1 {
			err = fmt.Errorf("unable to add a record, not able to find the object types for object %s, please try again later or contact the support", *apiName)
			return
		}
		result.ObjectTypes = objectTypes
	} else {
		result.ObjectTypes = objectTypeProvided
	}
	return
}

func GetObjectInfoByRecordIdentifier(userId *int64, identifier *string) (result *structs.Entity, err error) {
	decoded, errB := base64.StdEncoding.DecodeString(*identifier)
	if errB == nil {
		identiferParts := strings.Split(string(decoded), "|")
		if len(identiferParts) == 4 {
			objectId, errP := strconv.ParseInt(identiferParts[1], 10, 64)
			if errP == nil {
				object, errO := GetEntity(userId, &objectId)
				if errO == nil {
					result = object
				} else {
					log.Print(errO.Error(), errO)
					err = fmt.Errorf("unable to locate object information with identifer %s", *identifier)
				}
			} else {
				log.Print(errP.Error(), errP)
				err = fmt.Errorf("unable to identify object information with wrongly formatted identifer %s", *identifier)
			}
		} else {
			err = fmt.Errorf("unable to identify object information with wrongly formatted identifer %s", *identifier)
		}
	} else {
		log.Print(errB.Error(), errB)
		err = fmt.Errorf("unable to identify object information with wrongly formatted identifer %s", *identifier)
	}
	return
}

func RemoveRecord(userId *int64, apiName *string, identifier *string, id *int64) (result *structs.Record, err error) {
	if apiName == nil {
		object, errO := GetObjectInfoByRecordIdentifier(userId, identifier)
		if errO != nil {
			log.Print(errO.Error(), errO)
			err = fmt.Errorf("unable to remove a record for object api name %s and record with wrongly formatted identifer %s", ampstrings.EmptyIfNil(apiName), ampstrings.EmptyIfNil(identifier))
			return
		}
		apiName = object.ApiName
	}
	recordInfo := structs.RecordInfo{
		Identifier: identifier,
	}
	success, errP := recordInfo.Parse()
	if errP == nil && success {
		res, errRR := database.RemoveRecord(userId, apiName, recordInfo.Id, nil)
		if errRR != nil {
			log.Println(errRR.Error(), errRR)
			if strings.Contains(errRR.Error(), "[1451]") {
				err = fmt.Errorf("unable to remove a record for object %s and record identifer %s or id %s, another record has a reference on it", *apiName, ampstrings.EmptyIfNil(identifier), ampstrings.EmptyIfNilInt64(recordInfo.Id))
			} else {
				err = fmt.Errorf("unable to remove a record for object %s and record identifer %s or id %s", *apiName, ampstrings.EmptyIfNil(identifier), ampstrings.EmptyIfNilInt64(recordInfo.Id))
			}
		} else {
			result = res
		}
	} else {
		util.Loggify(errP)
		err = fmt.Errorf("supplied identifying information is not valid")
	}
	return result, err
}

func UpdateRecord(userId *int64, apiName *string, payload *string) (result *structs.Record, err error) {
	payloadData, errJ := jsons.GetJsonObject(payload)
	if errJ != nil {
		util.Loggify(errJ)
		err = fmt.Errorf("unable to update a record, not able to parse the supplied payload, please fix and try again or contact the support")
		return
	}
	identifier := payloadData[constants.GetDefaultFields().IDENTIFIER.ApiName]
	if identifier == nil || len(identifier.(string)) < 1 {
		err = fmt.Errorf("unable to update a record, the supplied payload does not contain record identifier, please fix and try again or contact the support")
		return
	}
	var object *structs.Entity
	var errO error
	if apiName == nil {
		object, errO = GetObjectInfoByRecordIdentifier(userId, util.PointerString(identifier.(string)))
		if errO != nil {
			log.Print(errO.Error(), errO)
			err = fmt.Errorf("unable to update a record for object api name %s and record identifer %s", *apiName, identifier.(string))
			return
		}
		apiName = object.ApiName
	}
	objectMetadata, errOM := GetObjectMetadata(userId, apiName, object, nil, nil)
	if errOM != nil {
		util.Loggify(errO)
		err = fmt.Errorf("unable to update a record for object api name %s and record identifer %s, not a valid identifier", *apiName, identifier.(string))
		return
	}

	recordInfo := structs.RecordInfo{
		Identifier: util.PointerString(identifier.(string)),
	}
	recordInfo.Parse()
	data, errR := database.FetchRecord(userId, apiName, recordInfo.Id, nil)
	if errR != nil || data.ID < 1 {
		util.Loggify(errR)
		err = fmt.Errorf("unable to update a record, the supplied payload does not contain a valid record identifier, please fix and try again or contact the support")
		return
	}
	objectTypeKey := constants.GetBaseObjectType().KEY
	if data.Record[constants.GetDefaultFields().OBJECTTYPE.ApiName] != nil {
		if reflect.TypeOf(data.Record[constants.GetDefaultFields().OBJECTTYPE.ApiName]).String() != "string" {
			err = fmt.Errorf("unable to add a record, the supplied object type info is wrong, please fix and try again or contact the support")
			return
		} else {
			objectTypeKey = data.Record[constants.GetDefaultFields().OBJECTTYPE.ApiName].(string)
		}
	}
	var objectType *structs.ObjectType
	for _, value := range *objectMetadata.ObjectTypes {
		if *value.ApiName == objectTypeKey {
			objectType = &value
			break
		}
	}
	if objectType == nil {
		err = fmt.Errorf("unable to update a record, not able to locate the object type info, please fix and try again or contact the support")
		return
	}

	fieldsMap := fieldsToMap(objectMetadata.Fields)
	var updateRecord = make(map[string]string)
	var errorValues = make([]string, 0)
	for _, value := range objectType.ObjectTypeFields {
		field := fieldsMap[*value.FieldId]
		if field.Status == nil || *field.Status != 1 || *value.ApiName == constants.GetDefaultFields().ID.ApiName || *value.ApiName == constants.GetDefaultFields().IDENTIFIER.ApiName || *value.ApiName == constants.GetDefaultFields().OBJECTTYPE.ApiName {
			continue
		}
		dirtyValue, okV := payloadData[*field.ApiName]
		if okV {
			var convertedValue string = ""
			var errC error
			if dirtyValue != nil {
				convertedValue, errC = ConvertFieldValue(&dirtyValue, &field)
				if errC != nil {
					errorValues = append(errorValues, errC.Error())
				}
			}
			if *field.Required > 0 && len(convertedValue) < 1 {
				errorValues = append(errorValues, fmt.Sprintf("The %s is a required field, it must be supplied with a value of type %s", *field.ApiName, *field.Type))
			} else {
				updateRecord[*field.ApiName] = convertedValue
				data.Record[*field.ApiName] = convertedValue
			}

		}
	}
	if len(errorValues) > 0 {
		return nil, fmt.Errorf(strings.Join(errorValues, ". "))
	} else if len(updateRecord) < 1 {
		return nil, fmt.Errorf("unable to update a record, no valid fields were supplied for update")
	}
	updateRecord[constants.GetDefaultFields().IDENTIFIER.ApiName] = identifier.(string)
	success, errU := database.UpdateRecord(userId, apiName, &updateRecord)
	if success {
		result = &data
	} else {
		util.Loggify(errU)
		err = fmt.Errorf("unable to update a record, not able to exute the updated changes, please fix and try again or contact the support")
	}
	return result, err
}

func AddRecords(userId *int64, apiName *string, payloads *string) (resultSuccess *[]map[string]interface{}, resultError *[]map[string]interface{}, err error) {
	objectMetadata, err := GetObjectMetadata(userId, apiName, nil, nil, nil)
	if err != nil {
		return
	}
	payloadsData, errJ := jsons.GetJsonArray(payloads)
	if errJ != nil {
		err = fmt.Errorf("unable to add records, not able to parse array from the supplied payload, please fix and try again or contact the support")
		return
	}

	resultSuccessList := list.New()
	resultErrorList := list.New()

	objectTypePartitions := make(map[int64]*list.List)
	fieldsMap := fieldsToMap(objectMetadata.Fields)
	objectTypeMap := objectTypesToMap(objectMetadata.ObjectTypes)
	for _, payloadData := range payloadsData {
		objectTypeKey := constants.GetBaseObjectType().KEY
		if payloadData[constants.GetDefaultFields().OBJECTTYPE.ApiName] != nil {
			if reflect.TypeOf(payloadData[constants.GetDefaultFields().OBJECTTYPE.ApiName]).String() != "string" {
				payloadData["error"] = "the supplied object type is of wrong data type"
				resultErrorList.PushBack(structs.Record{
					Record: payloadData,
				})
				continue
			} else {
				objectTypeKey = payloadData[constants.GetDefaultFields().OBJECTTYPE.ApiName].(string)
			}
		}
		objectType := objectTypeMap[objectTypeKey]

		if objectType == nil {
			payloadData["error"] = "the supplied object type is not valid"
			resultErrorList.PushBack(structs.Record{
				Record: payloadData,
			})
			continue
		}
		partitions := objectTypePartitions[*objectType.ID]
		var partition *list.List
		if partitions == nil {
			partitions = list.New()
			partition = list.New()
			partitions.PushBack(partition)
			objectTypePartitions[*objectType.ID] = partitions
		} else {
			for partEl := partitions.Front(); partEl != nil; partEl = partEl.Next() {
				if partEl.Value != nil {
					part := partEl.Value.(*list.List)
					if part.Len() <= 300 {
						partition = part
						break
					}
				}
			}
			if partition == nil {
				partition = list.New()
				partitions.PushBack(partition)
			}
		}

		var createRecord = make(map[string]string)
		for _, value := range objectType.ObjectTypeFields {
			field := fieldsMap[*value.FieldId]
			if field.Status == nil || *field.Status != 1 || *value.ApiName == constants.GetDefaultFields().ID.ApiName || *value.ApiName == constants.GetDefaultFields().IDENTIFIER.ApiName || *value.ApiName == constants.GetDefaultFields().OBJECTTYPE.ApiName {
				continue
			}
			if payloadData[*field.ApiName] == nil {
				if *field.Required > 0 {
					errorMessage := fmt.Sprintf("the %s is a required field, it must be supplied with a value of type %s", *field.ApiName, *field.Type)
					if payloadData["error"] != nil {
						payloadData["error"] = payloadData["error"].(string) + ", " + errorMessage
					} else {
						payloadData["error"] = errorMessage
					}
				}
			} else {
				dirtyValue := payloadData[*field.ApiName]
				convertedValue, errC := ConvertFieldValue(&dirtyValue, &field)
				if errC != nil {
					if *field.Required > 0 {
						if payloadData["error"] != nil {
							payloadData["error"] = payloadData["error"].(string) + ", " + errC.Error()
						} else {
							payloadData["error"] = errC.Error()
						}
					}
				} else {
					createRecord[*field.ApiName] = convertedValue
				}
			}
		}
		createRecord[constants.GetDefaultFields().OBJECTTYPE.ApiName] = *objectType.ApiName
		if payloadData["error"] == nil {
			createRecord[constants.GetDefaultFields().IDENTIFIER.ApiName] = strconv.FormatInt(*business.AmperId(), 10) + "|" + strconv.FormatInt(*objectMetadata.Object.ID, 10) + "|" + strconv.FormatInt(*objectType.ID, 10) + "|" + *util.UUID()
			partition.PushBack(createRecord)
		} else {
			resultErrorList.PushBack(structs.Record{
				Record: payloadData,
			})
		}
	}

	for _, partitions := range objectTypePartitions {
		if partitions != nil {
			for partition := partitions.Front(); partition != nil; partition = partition.Next() {
				if partition != nil {
					payloadPartition := partition.Value.(*list.List)
					resultRecords, errA := database.AddRecords(userId, apiName, payloadPartition)
					if errA != nil {
						for payload := payloadPartition.Front(); payload != nil; payload = payload.Next() {
							payloadData := payload.Value.(map[string]string)
							payloadData["error"] = "record addition failed due to database communication error and invalid formatted data"
							errorRecord := maps.GetStringToInterfaceMap(&payloadData)
							errorResult := structs.Record{
								Record: *errorRecord,
							}
							resultErrorList.PushBack(errorResult)
						}
					} else {
						partitionIds := make([]int64, resultRecords.Len())
						index := 0
						for resultRecord := resultRecords.Front(); resultRecord != nil; resultRecord = resultRecord.Next() {
							resultRecordData := resultRecord.Value.(structs.Record)
							partitionIds[index] = resultRecordData.ID
							index++
						}
						success, errU := database.ComputeRecordsIdentifier(userId, apiName, &partitionIds)
						if errU != nil && !success {
							for resultRecord := resultRecords.Front(); resultRecord != nil; resultRecord = resultRecord.Next() {
								resultRecordData := resultRecord.Value.(structs.Record)
								resultRecordData.Record["error"] = "not able to compute the identifier"
								resultErrorList.PushBack(resultRecordData.Record)
							}
							util.Loggify(errU)
							_, errR := RemoveRecords(userId, apiName, &partitionIds, nil)
							if errR != nil {
								util.Loggify(errR)
							}
						} else {
							for resultRecord := resultRecords.Front(); resultRecord != nil; resultRecord = resultRecord.Next() {
								resultRecordData := resultRecord.Value.(structs.Record)
								resultSuccessList.PushBack(resultRecordData)
							}
						}
					}
				}
			}
		}
	}

	temp := make([]map[string]interface{}, resultSuccessList.Len())
	index := 0
	for resultRecord := resultSuccessList.Front(); resultRecord != nil; resultRecord = resultRecord.Next() {
		resultRecordData := resultRecord.Value.(structs.Record)
		objectTypeKey := resultRecordData.Record[constants.GetDefaultFields().OBJECTTYPE.ApiName].(string)
		objectType := objectTypeMap[objectTypeKey]
		resultRecordData.Record[constants.GetDefaultFields().IDENTIFIER.ApiName] = GenerateRecordIdentifier(objectMetadata.Object.ID, objectType.ID, &resultRecordData.ID)
		temp[index] = resultRecordData.Record
		index++
	}
	resultSuccess = &temp

	temp1 := make([]map[string]interface{}, resultErrorList.Len())
	index = 0
	for resultRecord := resultErrorList.Front(); resultRecord != nil; resultRecord = resultRecord.Next() {
		resultRecordData := resultRecord.Value.(structs.Record)
		temp1[index] = resultRecordData.Record
		index++
	}
	resultError = &temp1
	return resultSuccess, resultError, err
}

func GenerateRecordIdentifier(objectId *int64, objectTypeId *int64, recordId *int64) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(*business.AmperId(), 10) + "|" + strconv.FormatInt(*objectId, 10) + "|" + strconv.FormatInt(*objectTypeId, 10) + "|" + strconv.FormatInt(*recordId, 10)))
}

func RemoveRecords(userId *int64, apiName *string, ids *[]int64, identifiers *[]string) (result bool, err error) {
	objectRecordIdsToRemove := make(map[string]*list.List)
	if apiName != nil && ids != nil {
		partitionList := list.New()
		partition := list.New()
		partitionList.PushBack(partition)
		for _, value := range *ids {
			if partition.Len() > 500 {
				partition = list.New()
				partitionList.PushBack(partition)
			}
			partition.PushBack(value)
		}
		objectRecordIdsToRemove[*apiName] = partitionList
	}
	if identifiers != nil {
		objectApiName := make(map[int64]*string)
		for _, identifier := range *identifiers {
			recordInfo := structs.RecordInfo{
				Identifier: &identifier,
			}
			success, errI := recordInfo.Parse()
			if errI != nil {
				util.Loggify(errI)
			} else if success {
				if recordInfo.ObjectId != nil {
					if objectApiName[*recordInfo.ObjectId] == nil {
						object, errO := GetEntity(userId, recordInfo.ObjectId)
						if err != nil || object.ApiName == nil {
							util.Loggify(errO)
							continue
						} else {
							objectApiName[*recordInfo.ObjectId] = object.ApiName
						}
					}
					apiName := objectApiName[*recordInfo.ObjectId]
					if recordInfo.Id != nil && apiName != nil {
						partitionList := objectRecordIdsToRemove[*apiName]
						var partition *list.List
						if partitionList == nil {
							partitionList = list.New()
							partition = list.New()
							partitionList.PushBack(partition)
							objectRecordIdsToRemove[*apiName] = partitionList
						} else {
							for currentPartition := partitionList.Front(); currentPartition != nil; currentPartition = currentPartition.Next() {
								if (currentPartition.Value.(*list.List)).Len() < 500 {
									partition = currentPartition.Value.(*list.List)
								}
							}
							if partition == nil {
								partition = list.New()
								partitionList.PushBack(partition)
							}
						}
						partition.PushBack(recordInfo.Id)
					}
				}
			}
		}
	}
	var errors strings.Builder
	for apiNameKey, partitions := range objectRecordIdsToRemove {
		for partition := partitions.Front(); partition != nil; partition = partition.Next() {
			success, errR := database.RemoveRecords(userId, &apiNameKey, partition.Value.(*list.List))
			if errR != nil || !success {
				result = false
				ids := make([]string, partition.Value.(*list.List).Len())
				index := 0
				for item := partition.Value.(*list.List).Front(); item != nil; item = item.Next() {
					ids[index] = strconv.FormatInt(*item.Value.(*int64), 10)
					index++
				}
				if strings.Contains(errR.Error(), "[1451]") {
					errors.WriteString(fmt.Sprintf("unable to remove %s records with ids : %s because one or some of the records have reference from another object record", apiNameKey, strings.Join(ids, ",")))
				} else {
					errors.WriteString(fmt.Sprintf("unable to remove %s records with ids : %s", apiNameKey, strings.Join(ids, ", ")))
				}
			}
		}
	}
	if errors.Len() > 0 {
		err = fmt.Errorf(errors.String())
	} else {
		result = true
	}
	return result, err
}

func UpdateRecords(userId *int64, payloads *string) (resultSuccess []map[string]interface{}, resultError []map[string]interface{}, err error) {
	payloadData, errJ := jsons.GetJsonArray(payloads)
	if errJ != nil || payloadData == nil {
		util.Loggify(errJ)
		err = fmt.Errorf("unable to update records, not able to parse the supplied payloads, not a valid json array, please fix and try again or contact the support")
		return
	}
	objectMetadataCache := make(map[string]*structs.ObjectMetadata)
	existingRecordIdsCache := make(map[string]*list.List)
	payloadRecordsCache := make(map[string]*map[int64]*map[string]interface{})
	existingRecordsCache := make(map[string]*map[int64]*structs.Record)

	for l := 0; l < len(payloadData); l++ {
		payload := payloadData[l]
		identifier := payload[constants.GetDefaultFields().IDENTIFIER.ApiName]
		if identifier == nil || len(identifier.(string)) < 1 {
			payload["error"] = "unable to update a record, the supplied payload does not contain record identifier, please fix and try again or contact the support"
			resultError = append(resultError, payload)
			continue
		}
		object, errO := GetObjectInfoByRecordIdentifier(userId, util.PointerString(identifier.(string)))
		if errO != nil {
			log.Print(errO.Error(), errO)
			payload["error"] = "unable to update a record, the supplied payload does not contain valid record identifier, please fix and try again or contact the support"
			resultError = append(resultError, payload)
			continue
		}
		var objectMetadata = objectMetadataCache[*object.ApiName]
		if objectMetadata == nil {
			tempObjectMetadata, errOM := GetObjectMetadata(userId, object.ApiName, object, nil, nil)
			if errOM != nil {
				util.Loggify(errO)
				payload["error"] = "unable to update a record, the supplied payload does not contain valid record identifier, please fix and try again or contact the support"
				resultError = append(resultError, payload)
				continue
			}
			objectMetadataCache[*object.ApiName] = &tempObjectMetadata
			objectMetadata = &tempObjectMetadata
		}
		recordInfo := structs.RecordInfo{
			Identifier: util.PointerString(identifier.(string)),
		}
		recordInfo.Parse()
		recordIds := existingRecordIdsCache[*object.ApiName]
		if recordIds == nil {
			recordIds = list.New()
			existingRecordIdsCache[*object.ApiName] = recordIds
		}
		recordIds.PushBack(recordInfo.Id)
		payloadRecords := payloadRecordsCache[*object.ApiName]
		if payloadRecords == nil {
			payloadRecordsTemp := make(map[int64]*map[string]interface{})
			payloadRecordsCache[*object.ApiName] = &payloadRecordsTemp
			payloadRecords = &payloadRecordsTemp
		}
		(*payloadRecords)[*recordInfo.Id] = &payload
	}
	for objectApiName, recordIds := range existingRecordIdsCache {
		inIDs := ampstrings.JoinListInt64(recordIds, ",")
		searchParams := make(map[string]interface{})
		inMap := make(map[string]interface{})
		inMap["id"] = *inIDs
		searchParams["in"] = inMap
		startId := int64(0)
		limit := int64(recordIds.Len())
		existingRecords, _, errER := database.FetchRecords(userId, &objectApiName, &startId, &limit, searchParams, nil, nil)
		if errER != nil || existingRecords == nil {
			util.Loggify(errER)
			payloadRecords := payloadRecordsCache[objectApiName]
			if payloadRecords != nil {
				for recordId := recordIds.Front(); recordId != nil; recordId = recordId.Next() {
					payload := (*payloadRecords)[*recordId.Value.(*int64)]
					if payload != nil {
						(*payload)["error"] = "unable to update a record, the supplied payload is not a valid record, please fix and try again or contact the support"
						resultError = append(resultError, *payload)
					}
				}
			}
			delete(payloadRecordsCache, objectApiName)
		} else {
			existingRecordsMap := existingRecordsCache[objectApiName]
			if existingRecordsMap == nil {
				existingRecordsMapTemp := make(map[int64]*structs.Record)
				existingRecordsCache[objectApiName] = &existingRecordsMapTemp
				existingRecordsMap = &existingRecordsMapTemp
			}

			for i := 0; i < len(*existingRecords); i++ {
				existingRecord := (*existingRecords)[i]
				id, _ := strconv.ParseInt(existingRecord["id"].(string), 10, 64)
				(*existingRecordsMap)[id] = &structs.Record{
					ID:     id,
					Record: existingRecord,
				}
			}
			objectPayloadRecordsCache := payloadRecordsCache[objectApiName]
			for payloadRecordId, payloadRecord := range *objectPayloadRecordsCache {
				existingRecord := (*existingRecordsMap)[payloadRecordId]
				if existingRecord == nil {
					if payloadRecord != nil {
						(*payloadRecord)["error"] = "unable to update a record, the supplied record does not exist, please fix and try again or contact the support"
						resultError = append(resultError, *payloadRecord)
					}
				}
			}
		}
	}
	recordsToUpdateMap := make(map[string]*list.List)
	for objectApiName, payloadRecords := range payloadRecordsCache {
		objectMetadata := objectMetadataCache[objectApiName]
		objectTypeMap := objectTypesToMap(objectMetadata.ObjectTypes)

		existingRecords := existingRecordsCache[objectApiName]
		hasError := false
		for recordId, payloadRecord := range *payloadRecords {
			recordToUpdate := make(map[string]string)
			existingRecord := (*existingRecords)[recordId]
			if existingRecord != nil && existingRecord.Record["id"] != nil {
				var objectType *structs.ObjectType
				objectTypeName := existingRecord.Record[constants.GetDefaultFields().OBJECTTYPE.ApiName]
				if objectTypeName != nil && len(objectTypeName.(string)) > 0 {
					objectType = objectTypeMap[objectTypeName.(string)]
				}
				if objectType != nil {
					fieldsMap := fieldsToMap(objectMetadata.Fields)
					for _, objectTypeField := range objectType.ObjectTypeFields {
						field := fieldsMap[*objectTypeField.FieldId]
						if field.Status == nil || *field.Status != int8(1) || *objectTypeField.ApiName == constants.GetDefaultFields().ID.ApiName || *objectTypeField.ApiName == constants.GetDefaultFields().IDENTIFIER.ApiName || *objectTypeField.ApiName == constants.GetDefaultFields().OBJECTTYPE.ApiName {
							continue
						}
						dirtyValue, isDirtyValueOk := (*payloadRecord)[*field.ApiName]
						if isDirtyValueOk {
							if dirtyValue == nil && *field.Required > 0 {
								if payloadRecord != nil {
									(*payloadRecord)["error"] = fmt.Sprintf("unable to update a record, the supplied payload record has has a required field %s which is empty, please fix and try again or contact the support", *field.ApiName)
									resultError = append(resultError, *payloadRecord)
								}
								hasError = true
								break
							} else {
								convertedValue, errC := ConvertFieldValue(&dirtyValue, &field)
								if errC != nil {
									if payloadRecord != nil {
										(*payloadRecord)["error"] = errC.Error()
										resultError = append(resultError, *payloadRecord)
									}
									hasError = true
									break
								} else if len(convertedValue) < 1 && *field.Required == int8(1) {
									if payloadRecord != nil {
										(*payloadRecord)["error"] = fmt.Sprintf("unable to update a record, the supplied payload record has has a required field %s which is empty, please fix and try again or contact the support", *field.ApiName)
										resultError = append(resultError, *payloadRecord)
									}
									hasError = true
									break
								} else {
									recordToUpdate[*field.ApiName] = convertedValue
								}
							}
						}
					}
				} else {
					if payloadRecord != nil {
						(*payloadRecord)["error"] = "unable to update a record, the supplied payload record has no valid object type, please fix and try again or contact the support"
						resultError = append(resultError, *payloadRecord)
					}
					hasError = true
				}
			} else {
				if payloadRecord != nil {
					(*payloadRecord)["error"] = "unable to update a record, the supplied payload does not exist, please fix and try again or contact the support"
					resultError = append(resultError, *payloadRecord)
				}
				hasError = true
			}
			if !hasError {
				recordsToUpdatePartions := recordsToUpdateMap[objectApiName]
				var recordsToUpdatePartion *list.List
				if recordsToUpdatePartions == nil {
					recordsToUpdatePartions = list.New()
					recordsToUpdateMap[objectApiName] = recordsToUpdatePartions
				}
				for recordsToUpdatePartionEl := recordsToUpdatePartions.Front(); recordsToUpdatePartionEl != nil; recordsToUpdatePartionEl = recordsToUpdatePartionEl.Next() {
					currentRecordsToUpdatePartion := recordsToUpdatePartionEl.Value.(*list.List)
					if currentRecordsToUpdatePartion == nil || currentRecordsToUpdatePartion.Len() < 300 {
						recordsToUpdatePartion = currentRecordsToUpdatePartion
					}
				}
				if recordsToUpdatePartion == nil {
					recordsToUpdatePartion = new(list.List)
					recordsToUpdatePartions.PushBack(recordsToUpdatePartion)
				}
				recordToUpdate["id"] = strconv.FormatInt(recordId, 10)
				recordsToUpdatePartion.PushBack(&recordToUpdate)
			}
		}
	}
	for objectApiName, recordsToUpdatePartions := range recordsToUpdateMap {
		for recordsToUpdatePartion := recordsToUpdatePartions.Front(); recordsToUpdatePartion != nil; recordsToUpdatePartion = recordsToUpdatePartion.Next() {
			success, errP := database.UpdateRecords(userId, &objectApiName, recordsToUpdatePartion.Value.(*list.List))
			if errP != nil || !success {
				util.Loggify(errP)
				for recordsToUpdateEl := recordsToUpdatePartion.Value.(*list.List).Front(); recordsToUpdateEl != nil; recordsToUpdateEl = recordsToUpdateEl.Next() {
					erroredRecord := recordsToUpdateEl.Value.(*map[string]string)
					if erroredRecord != nil {
						erroredRecordInt := *maps.GetStringToInterfaceMap(erroredRecord)
						erroredRecordInt["error"] = "unable to update a record, there was a update error, please fix and try again or contact the support"
						resultError = append(resultError, erroredRecordInt)
					}
				}
			} else {
				if recordsToUpdatePartion.Value != nil {
					for recordsToUpdateEl := recordsToUpdatePartion.Value.(*list.List).Front(); recordsToUpdateEl != nil; recordsToUpdateEl = recordsToUpdateEl.Next() {
						succeedRecord := recordsToUpdateEl.Value.(*map[string]string)
						if succeedRecord != nil {
							id, _ := strconv.ParseInt((*succeedRecord)["id"], 10, 64)
							succeedRecordInt := *maps.GetStringToInterfaceMap(succeedRecord)
							successUpdate := structs.Record{
								ID:     id,
								Record: succeedRecordInt,
							}
							existingRecords := existingRecordsCache[objectApiName]
							if existingRecords != nil {
								existingRecord := (*existingRecords)[id]
								if existingRecord != nil && existingRecord.Record != nil {
									for key, value := range succeedRecordInt {
										existingRecord.Record[key] = value
									}
									successUpdate.Record = existingRecord.Record
								}
							}
							resultSuccess = append(resultSuccess, successUpdate.Record)
						}
					}
				}
			}
		}
	}

	return resultSuccess, resultError, err
}
