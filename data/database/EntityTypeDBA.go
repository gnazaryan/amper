package database

import (
	databasecache "amper/cache/database"
	"amper/common/constants"
	"amper/common/structs"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

const GET_OBJECT_TYPE_BY_API_NAME = "SELECT * FROM amper.object_type_sys as OT where object_id=%d and apiName='%s'"

func GetObjectTypeByApiName(objectId *int64, apiName *string) (*structs.ObjectType, error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(GET_OBJECT_TYPE_BY_API_NAME, *objectId, *apiName)
	rows, errQ := pool.Query(query)
	if errQ == nil {
		for rows.Next() {
			var result structs.ObjectType
			rows.Scan(&result.ID, &result.ObjectId, &result.ApiName, &result.Label, &result.ExtendsTo, &result.CreatedBy, &result.CreatedDate)
			rows.Close()
			return &result, nil
		}
	} else {
		var err = fmt.Errorf("unable to run query against database to get object type with api name: %s", *apiName)
		log.Print(errQ.Error(), errQ)
		return nil, err
	}
	return nil, nil
}

const GET_OBJECT_TYPE_BY_ID = "SELECT * FROM amper.object_type_sys WHERE id=%d"

func GetObjectType(userID *int64, objectTypeId *int64) (*structs.ObjectType, error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(GET_OBJECT_TYPE_BY_ID, *objectTypeId)
	rows, errQ := pool.Query(query)
	if errQ == nil {
		for rows.Next() {
			var result structs.ObjectType
			rows.Scan(&result.ID, &result.ObjectId, &result.ApiName, &result.Label, &result.ExtendsTo, &result.CreatedBy, &result.CreatedDate)
			rows.Close()
			return &result, nil
		}
	} else {
		log.Print(errQ.Error(), errQ)
	}
	return nil, fmt.Errorf("unable to run query against database to get object type with id: %d", *objectTypeId)
}

const GET_OBJECT_TYPES_BY_OBJECTID = "SELECT * FROM amper.object_type_sys WHERE object_id=%d"

func GetObjectTypes(userID *int64, objectId *int64) (result []structs.ObjectType, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(GET_OBJECT_TYPES_BY_OBJECTID, *objectId)
	rows, errQ := pool.Query(query)
	if errQ == nil {
		for rows.Next() {
			var objectType structs.ObjectType
			rows.Scan(&objectType.ID, &objectType.ObjectId, &objectType.ApiName, &objectType.Label, &objectType.ExtendsTo, &objectType.CreatedBy, &objectType.CreatedDate)
			result = append(result, objectType)
		}
	} else {
		log.Print(errQ.Error(), errQ)
		err = fmt.Errorf("uneble to retrieve the object types information")
	}
	return
}

// GetEntities query the database and retrieve the entities for the provided start and limit parameters
func GetObjectTypesFields(userID *int64, objectId *int64) (result *[]structs.ObjectType, err error) {
	objectTypeFields, errDB := GetAllObjectTypeFields(userID, objectId)
	if errDB == nil {
		result = new([]structs.ObjectType)
		for _, value := range objectTypeFields {
			var objectType structs.ObjectType
			var objectTypeInfo = value[0]
			objectType.ID = objectTypeInfo.ObjectTypeId
			objectType.ApiName = objectTypeInfo.ObjectTypeApiName
			objectType.Label = objectTypeInfo.ObjectTypeLabel

			//lookup for the extends to object type label, it must exist
			if objectTypeInfo.ObjectTypeExtendsTo != nil {
				extendsToObjectType := objectTypeFields[*objectTypeInfo.ObjectTypeExtendsTo]
				if len(extendsToObjectType) > 0 {
					objectType.ExtendsToLabel = extendsToObjectType[0].ObjectTypeLabel
				}
			}
			objectType.ExtendsTo = objectTypeInfo.ObjectTypeExtendsTo
			objectType.ObjectId = objectId
			objectType.CreatedBy = objectTypeInfo.ObjectTypeCreatedBy
			objectType.CreatedDate = objectTypeInfo.ObjectTypeCreatedDate
			var resultObjectTypeFields = make([]structs.ObjectTypeField, 0)
			constructObjectTypeFields(&objectTypeFields, &value, &resultObjectTypeFields)
			objectType.ObjectTypeFields = resultObjectTypeFields
			*result = append(*result, objectType)
		}
	} else {
		err = errors.New("unable to run query against database to get object types and object type fields")
		log.Print(errDB.Error(), err)
	}
	return
}

func GetObjectTypeFields(userID *int64, objectId *int64, objectTypeId *int64) (result *structs.ObjectType, err error) {
	objectTypeFields, errDB := GetObjectTypeFieldsByObjectTypeId(userID, objectId, objectTypeId)
	if errDB == nil {
		for _, value := range objectTypeFields {
			var objectType structs.ObjectType
			var objectTypeInfo = value[0]
			objectType.ID = objectTypeInfo.ObjectTypeId
			objectType.ApiName = objectTypeInfo.ObjectTypeApiName
			objectType.Label = objectTypeInfo.ObjectTypeLabel

			//lookup for the extends to object type label, it must exist
			if objectTypeInfo.ObjectTypeExtendsTo != nil {
				extendsToObjectType := objectTypeFields[*objectTypeInfo.ObjectTypeExtendsTo]
				if len(extendsToObjectType) > 0 {
					objectType.ExtendsToLabel = extendsToObjectType[0].ObjectTypeLabel
				}
			}
			objectType.ExtendsTo = objectTypeInfo.ObjectTypeExtendsTo
			objectType.ObjectId = objectId
			objectType.CreatedBy = objectTypeInfo.ObjectTypeCreatedBy
			objectType.CreatedDate = objectTypeInfo.ObjectTypeCreatedDate
			var resultObjectTypeFields = make([]structs.ObjectTypeField, 0)
			constructObjectTypeFields(&objectTypeFields, &value, &resultObjectTypeFields)
			objectType.ObjectTypeFields = resultObjectTypeFields
			result = &objectType
		}
	} else {
		err = errors.New("unable to run query against database to get object types and object type fields")
		log.Print(errDB.Error(), err)
	}
	/*var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(GET_OBJECT_TYPES_FOR_OBJECT, *objectId)
	rows, errQ := pool.Query(query)
	if errQ == nil {
		for rows.Next() {
			var objectType structs.ObjectType
			rows.Scan(&objectType.ID, &objectType.ObjectId, &objectType.ApiName, &objectType.Label, &objectType.ExtendsTo, &objectType.CreatedBy, &objectType.CreatedDate)
			result = append(result, objectType)
		}
	} else {
		err = errors.New("unable to run query against database to get entities")
		log.Print(errQ.Error(), errQ)
	}*/
	return
}

func constructObjectTypeFields(objectTypeFields *map[int64][]structs.ObjectTypeField, values *[]structs.ObjectTypeField, result *[]structs.ObjectTypeField) {
	var objectTypeInfo = (*values)[0]
	if objectTypeInfo.ObjectTypeExtendsTo != nil {
		var extendsToValues = (*objectTypeFields)[*objectTypeInfo.ObjectTypeExtendsTo]
		constructObjectTypeFields(objectTypeFields, &extendsToValues, result)
	}

	for _, objectTypeFieldDB := range *values {
		var objectTypeField structs.ObjectTypeField
		if objectTypeFieldDB.ID != nil {
			objectTypeField.ID = objectTypeFieldDB.ID
			objectTypeField.Label = objectTypeFieldDB.Label
			objectTypeField.FieldId = objectTypeFieldDB.FieldId
			objectTypeField.ApiName = objectTypeFieldDB.ApiName
			objectTypeField.Type = objectTypeFieldDB.Type
			objectTypeField.CreatedBy = objectTypeFieldDB.CreatedBy
			objectTypeField.CreatedDate = objectTypeFieldDB.CreatedDate
			objectTypeField.ObjectTypeExtendsTo = objectTypeFieldDB.ObjectTypeExtendsTo
			objectTypeField.ObjectTypeId = objectTypeFieldDB.ObjectTypeId
			objectTypeField.ObjectTypeLabel = objectTypeFieldDB.ObjectTypeLabel
			*result = append(*result, objectTypeField)
		}
	}
}

const GET_OBJECT_TYPE_FIELDS_FOR_OBJECT = "SELECT OTF.id, OTF.field_id, OTF.created_by, OTF.created_date, F.label, F.apiName, F.type, OT.apiName, OT.label, OT.extends, OT.id, OT.created_by, OT.created_date FROM amper.object_type_field_sys as OTF INNER JOIN amper.field_sys as F on F.id=OTF.field_id RIGHT JOIN amper.object_type_sys as OT on OT.id=OTF.object_type_id where OT.object_id=%d"

func GetAllObjectTypeFields(userID *int64, objectId *int64) (result map[int64][]structs.ObjectTypeField, err error) {
	result = make(map[int64][]structs.ObjectTypeField)
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(GET_OBJECT_TYPE_FIELDS_FOR_OBJECT, *objectId)
	rows, errQ := pool.Query(query)
	if errQ == nil {
		for rows.Next() {
			var objectTypeField structs.ObjectTypeField
			rows.Scan(&objectTypeField.ID,
				&objectTypeField.FieldId, &objectTypeField.CreatedBy, &objectTypeField.CreatedDate, &objectTypeField.Label, &objectTypeField.ApiName, &objectTypeField.Type,
				&objectTypeField.ObjectTypeApiName, &objectTypeField.ObjectTypeLabel, &objectTypeField.ObjectTypeExtendsTo, &objectTypeField.ObjectTypeId,
				&objectTypeField.ObjectTypeCreatedBy, &objectTypeField.ObjectTypeCreatedDate)

			ComputeIfApsent(&result, objectTypeField.ObjectTypeId, &objectTypeField)
		}
	} else {
		err = errors.New("unable to run query against database to get object type fields")
		log.Print(errQ.Error(), errQ)
	}
	rows.Close()
	return
}

const GET_OBJECT_TYPE_FIELDS_FOR_OBJECT_TYPE = "SELECT OTF.id, OTF.field_id, OTF.created_by, OTF.created_date, F.label, F.apiName, F.type, OT.apiName, OT.label, OT.extends, OT.id, OT.created_by, OT.created_date FROM amper.object_type_field_sys as OTF INNER JOIN amper.field_sys as F on F.id=OTF.field_id RIGHT JOIN amper.object_type_sys as OT on OT.id=OTF.object_type_id where OT.object_id=%d AND OT.id=%d"

func GetObjectTypeFieldsByObjectTypeId(userID *int64, objectId *int64, objectTypeId *int64) (result map[int64][]structs.ObjectTypeField, err error) {
	result = make(map[int64][]structs.ObjectTypeField)
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(GET_OBJECT_TYPE_FIELDS_FOR_OBJECT_TYPE, *objectId, *objectTypeId)
	rows, errQ := pool.Query(query)
	if errQ == nil {
		for rows.Next() {
			var objectTypeField structs.ObjectTypeField
			rows.Scan(&objectTypeField.ID,
				&objectTypeField.FieldId, &objectTypeField.CreatedBy, &objectTypeField.CreatedDate, &objectTypeField.Label, &objectTypeField.ApiName, &objectTypeField.Type,
				&objectTypeField.ObjectTypeApiName, &objectTypeField.ObjectTypeLabel, &objectTypeField.ObjectTypeExtendsTo, &objectTypeField.ObjectTypeId,
				&objectTypeField.ObjectTypeCreatedBy, &objectTypeField.ObjectTypeCreatedDate)

			ComputeIfApsent(&result, objectTypeField.ObjectTypeId, &objectTypeField)
		}
	} else {
		err = errors.New("unable to run query against database to get object type fields")
		log.Print(errQ.Error(), errQ)
	}
	rows.Close()
	return
}

func ComputeIfApsent(input *map[int64][]structs.ObjectTypeField, key *int64, value *structs.ObjectTypeField) {
	temp := (*input)[*key]
	if temp == nil {
		var initialize []structs.ObjectTypeField
		initialize = append(initialize, *value)
		(*input)[*key] = initialize
	} else {
		temp = append(temp, *value)
		(*input)[*key] = temp
	}
}

const INSERT_OBJECT_TYPE_1 = "INSERT INTO amper.object_type_sys VALUES (null, %d, '%s', '%s', %d, %d, '%s')"

func CreateObjectType(UserID *int64, ObjectId *int64, ExtendsTo *int64, ApiName *string, Label *string) (result bool, err error) {
	timeFormatted := time.Now().Format(constants.TIME_FORMAT)
	insertObjectTypeStatement := fmt.Sprintf(INSERT_OBJECT_TYPE_1, *ObjectId, *ApiName, *Label, *ExtendsTo, *UserID, timeFormatted)
	var pool *sql.DB = databasecache.Pool()
	_, errDB := pool.Exec(insertObjectTypeStatement)
	if errDB != nil {
		result = false
		err = fmt.Errorf("unable to crete an object type for the api name: %s", *ApiName)
		log.Print(errDB.Error(), err)
	} else {
		result = true
	}
	return
}

const DELETE_OBJECT_TYPE_BY_ID = "DELETE from amper.object_type_sys where id='%d'"

func DeleteObjectType(UserID *int64, ObjectTypeId *int64) (result bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	if deleteObjectTypeFieldsByObjectTypeId(pool, ObjectTypeId) {
		_, errDB := pool.Exec(fmt.Sprintf(DELETE_OBJECT_TYPE_BY_ID, *ObjectTypeId))
		if errDB == nil {
			result = true
		} else {
			result = false
			err = fmt.Errorf("unable to delete an object type for the object type id: %d", *ObjectTypeId)
			log.Print(err.Error(), err)
		}
	} else {
		result = false
		err = fmt.Errorf("unable to delete an object type fields for the object type id: %d", *ObjectTypeId)
		log.Print(err.Error(), err)
	}
	return
}
