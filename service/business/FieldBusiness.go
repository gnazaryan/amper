package business

import (
	"amper/common/argument"
	"amper/common/constants"
	"amper/common/structs"
	"amper/common/util"
	"amper/common/util/arrays"
	"amper/data/database"
	"fmt"
	"log"
	"strconv"
)

func CreateField(userID *int64, entityId *int64, apiName *string, label *string, dataType *string, textLength *int64, objectReferenceId *int64, required *bool, status *bool) (bool, error) {
	var err error = nil
	if !argument.ValidateApiName(apiName) {
		err = fmt.Errorf("the provided api name: %s, is not valid, it should end with _amp", *apiName)
		return false, err
	}
	if constants.GetDataTypes.REFERENCE.Name == *dataType {
		err = argument.Validate(map[string]interface{}{"ObjectReferenceId": objectReferenceId})
	} else if constants.GetDataTypes.TEXT.Name == *dataType {
		err = argument.Validate(map[string]interface{}{"textLength": textLength})
	}
	if err != nil {
		return false, err
	}
	entity, errDb := database.GetEntity(entityId)
	if errDb != nil || entity.ID == nil {
		err = fmt.Errorf("unable to retrieve entity for the provided EntityId: %d", *entityId)
		log.Print(errDb.Error(), errDb)
		return false, err
	}
	var objectReferenceApiName *string
	if constants.GetDataTypes.REFERENCE.Name == *dataType {
		entityRference, errDb1 := database.GetEntity(objectReferenceId)
		objectReferenceApiName = entityRference.ApiName
		if errDb1 != nil || entityRference.ID == nil {
			err = fmt.Errorf("unable to retrieve entity for the provided ObjectReferenceId: %d", *objectReferenceId)
			log.Print(errDb.Error(), errDb)
			return false, err
		}
	}
	entityFields, errEf := database.GetEntityFields(entityId)
	if errEf != nil {
		err = fmt.Errorf("unable to retrieve entity fields for the provided EntityId: %d", *entityId)
		log.Print(errEf.Error(), errEf)
		return false, err
	}
	apiNames := entityFields.ApiNameList()
	validApiName := util.StringPointer(apiName)
	i := 1
	for arrays.Contains(apiNames, validApiName) {
		validApiName = util.AppendApiName(*apiName, strconv.Itoa(i))
		i++
	}

	return database.CreateEntityField(userID, entityId, entity.ApiName, &validApiName, label, dataType, textLength, objectReferenceId, objectReferenceApiName, required, status)
}

func DeleteField(userID *int64, entityId *int64, fieldIds *[]int64) (bool, error) {
	result, errEf := database.DeleteFields(entityId, fieldIds)
	if errEf != nil {
		log.Print(errEf.Error(), errEf)
		return false, errEf
	}
	return result, nil
}

func AddObjectTypeField(userID *int64, fieldId *int64, objectTypeId *int64) (bool, error) {
	objType, errOT := database.GetObjectType(userID, objectTypeId)
	if errOT != nil {
		log.Print(errOT.Error(), errOT)
		return false, fmt.Errorf("unable to add a field with id: %d to the specified object type: %d, because not able to identify if the object type is missing", *fieldId, *objectTypeId)
	}
	if objType == nil {
		return false, fmt.Errorf("unable to add a field with id: %d to the specified object type: %d, because the object type is missing", *fieldId, *objectTypeId)
	}
	objTypeFieldsMap, errorOTF := database.GetAllObjectTypeFields(userID, objType.ObjectId)
	if errorOTF != nil {
		return false, fmt.Errorf("unable to add a field with id: %d to the specified object type: %d, because not able to retrieve the object type hierarchy information", *fieldId, *objectTypeId)
	}
	if objTypeFieldsMap == nil {
		return false, fmt.Errorf("unable to add a field with id: %d to the specified object type: %d, because object type hierarchy information is missing", *fieldId, *objectTypeId)
	}

	objectTypesMap, errOTM := GetEntityTypesMap(userID, objType.ObjectId)
	if errOTM != nil {
		log.Print(errOT.Error(), errOT)
		return false, fmt.Errorf("unable to add a field with id: %d to the specified object type: %d, because not able to retrieve object types information", *fieldId, *objectTypeId)

	}

	//---------- Lookup the current object type for existing field, if field exist shouldn't let add again--
	objTypeFields := objTypeFieldsMap[*objType.ID]
	for _, objTypeField := range objTypeFields {
		if objTypeField.FieldId != nil && *objTypeField.FieldId == *fieldId {
			return false, fmt.Errorf("unable to add a field with id: %d to the specified object type: %d, because the object type already contains the field", *fieldId, *objectTypeId)
		}
	}
	//---------------------------------------------------------------------------------------------------
	//---------- Lookup upward the hierarchy for existing field, if field exist shouldn't let add again--
	var current *structs.ObjectType = objType
	var hierarchyObjTypes []int64
	for current != nil && current.ExtendsTo != nil && *current.ExtendsTo > 0 {
		hierarchyObjTypes = append(hierarchyObjTypes, *current.ExtendsTo)
		currentObjectType, contains := objectTypesMap[*current.ExtendsTo]
		if !contains {
			return false, fmt.Errorf("unable to add a field with id: %d to the specified object type: %d, because couldn't identify the hierarchy of object types", *fieldId, *objectTypeId)
		}
		current = &currentObjectType
	}
	for _, objTypeId := range hierarchyObjTypes {
		objTypeFields := objTypeFieldsMap[objTypeId]
		for _, objTypeField := range objTypeFields {
			if objTypeField.FieldId != nil && *objTypeField.FieldId == *fieldId {
				return false, fmt.Errorf("unable to add a field with id: %d to the specified object type: %d, because the UPWARD object type hierarchy already contains the field", *fieldId, *objectTypeId)
			}
		}
	}
	//---------------------------------------------------------------------------------------------------
	//---------- Lookup downward the hierarchy for existing field, if field exist shouldn't let add again--
	var downwardHierarchyObjTypes []int64
	var toVisitObjectTypes []structs.ObjectType
	toVisitObjectTypes = append(toVisitObjectTypes, *objType)
	for len(toVisitObjectTypes) > 0 {
		current, toVisitObjectTypes = RemoveFirst(toVisitObjectTypes)

		for objectTypeId, objectType := range objectTypesMap {
			if objectType.ExtendsTo != nil && *current.ID == *objectType.ExtendsTo {
				downwardHierarchyObjTypes = append(downwardHierarchyObjTypes, objectTypeId)
				toVisitObjectTypes = append(toVisitObjectTypes, objectType)
			}
		}
	}
	for _, objTypeId := range downwardHierarchyObjTypes {
		objTypeFields := objTypeFieldsMap[objTypeId]
		for _, objTypeField := range objTypeFields {
			if objTypeField.FieldId != nil && *objTypeField.FieldId == *fieldId {
				return false, fmt.Errorf("unable to add a field with id: %d to the specified object type: %d, because the DOWNWARD object type hierarchy already contains the field", *fieldId, *objectTypeId)
			}
		}
	}
	//---------------------------------------------------------------------------------------------------

	result, errEf := database.AddObjectTypeField(userID, fieldId, objectTypeId)
	if errEf != nil {
		err := fmt.Errorf("unable to add a field with id: %d to the specified object type: %d", *fieldId, *objectTypeId)
		log.Print(errEf.Error(), errEf)
		return false, err
	}
	return result, nil
}

func DeleteObjectTypeField(userID *int64, fieldId *int64, objectTypeId *int64) (bool, error) {
	result, errDOT := database.DeleteObjectTypeField(userID, fieldId, objectTypeId)
	if errDOT != nil {
		err := fmt.Errorf("unable to delete a field with id: %d to the specified object type: %d", *fieldId, *objectTypeId)
		log.Print(errDOT.Error(), errDOT)
		return false, err
	}
	return result, nil
}

func RemoveFirst(input []structs.ObjectType) (*structs.ObjectType, []structs.ObjectType) {
	var result structs.ObjectType
	var result1 []structs.ObjectType
	for index, item := range input {
		if index != 0 {
			result1 = append(result1, item)
		} else {
			result = item
		}
	}
	return &result, result1
}
