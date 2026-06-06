package database

import (
	databasecache "amper/cache/database"
	"amper/common/constants"
	"amper/common/structs"
	"amper/common/util"
	"amper/common/util/arrays"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

const deleteEntityFieldsByIds = "DELETE FROM amper.field_sys where entityId=%d and id in (%s)"
const dropForeginKeyOnTable = "ALTER TABLE `amper`.`%s` DROP FOREIGN KEY `%s`"
const getForeginKeyNames = "SELECT CONSTRAINT_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE REFERENCED_TABLE_SCHEMA = 'amper' AND TABLE_NAME = '%s' AND COLUMN_NAME = '%s';"
const dropColumnOnTable = "ALTER TABLE `amper`.`%s` DROP COLUMN `%s`"

func DeleteFields(entityId *int64, fieldIds *[]int64) (result bool, err error) {
	result = true
	var pool *sql.DB = databasecache.Pool()
	entity, errE := GetEntity(entityId)
	if errE != nil {
		log.Print(errE.Error(), errE)
		return false, errE
	} else if entity == nil || entity.ID == nil {
		errO := fmt.Errorf("unable to remove object field because no object found with the provided objectId: %d", *entityId)
		return false, errO
	}

	fields, errF := GetEntityFieldsByIds(entity.ID, fieldIds)
	if errF != nil {
		log.Print(errF.Error(), errF)
		return false, errF
	}

	for i := 0; i < len(fields); i++ {
		field := fields[i]
		apiName := field.ApiName
		dataType := field.Type
		if !constants.IsSystemField(apiName) {
			objTF, errOTF := GetObjectTypeFieldById(field.ID)
			if errOTF != nil || objTF != nil {
				err = fmt.Errorf("unable to remove an object field by api name: %s, because it has a reference from a object type, remove object type field first", *field.ApiName)
				result = false
				continue
			}

			if constants.GetDataTypes.REFERENCE.Name == *dataType {
				var constraintName string
				errC := pool.QueryRow(fmt.Sprintf(getForeginKeyNames, util.StringPointer(entity.ApiName), util.StringPointer(field.ApiName))).Scan(&constraintName)
				if errC != nil {
					errS := errors.New("unable to remove foregin constrain key, skipping field deletion, contact the adminstrator")
					log.Print(errS.Error(), errS)
					result = false
					continue
				}
				_, errF := pool.Exec(fmt.Sprintf(dropForeginKeyOnTable, util.StringPointer(entity.ApiName), constraintName))
				if errF != nil {
					errS := errors.New("unable to remove foregin constrain key, skipping field deletion, contact the adminstrator")
					log.Print(errS.Error(), errS)
					result = false
					continue
				}
			}
			_, errD := pool.Exec(fmt.Sprintf(dropColumnOnTable, util.StringPointer(entity.ApiName), util.StringPointer(field.ApiName)))
			if errD != nil {
				errD1 := errors.New("unable to delete a field from object")
				log.Print(errD1.Error(), errD1)
				log.Print(errD.Error(), errD)
			}
			res, errQ := pool.Exec(fmt.Sprintf(deleteEntityFieldsByIds, *entityId, strings.Join(arrays.IntToString(fieldIds), ",")))
			rowsAffected, errW := res.RowsAffected()
			if errQ == nil && errW == nil && rowsAffected == int64(len(*fieldIds)) {
				continue
			} else {
				errQ1 := errors.New("unable to remove a field from field__sys")
				log.Print(errQ1.Error(), errQ1)
				if errQ != nil {
					log.Print(errQ.Error(), errQ)
				}
				if errW != nil {
					log.Print(errW.Error(), errW)
				}
			}
		} else {
			errS := errors.New("unable to delete the entity field, skipping current field")
			log.Print(errS.Error(), errS)
			result = false
			continue
		}

		if !result {
			err = errors.New("some, or all of the fields were not able to be deleted, contact the support")
		}
		return result, err
	}
	return result, err
}

const getEntityFieldsByIds = "SELECT * FROM amper.field_sys where entityId=%d and id in (%s)"

func GetEntityFieldsByIds(entityId *int64, fieldIds *[]int64) (result structs.Fields, err error) {
	var pool *sql.DB = databasecache.Pool()
	rows, errQ := pool.Query(fmt.Sprintf(getEntityFieldsByIds, *entityId, strings.Join(arrays.IntToString(fieldIds), ",")))
	if errQ == nil {
		for rows.Next() {
			var field structs.Field
			rows.Scan(&field.ID, &field.ApiName, &field.Label, &field.Type, &field.Status, &field.Required, &field.ObjectId, &field.CreatedBy, &field.TextLength, &field.ObjectReference)
			result = append(result, field)
		}
	} else {
		err = errors.New("unable to run query against database to get entity fields")
		log.Print(errQ.Error(), errQ)
	}
	rows.Close()
	return
}

const getEntityFields = "SELECT * FROM amper.field_sys where entityId=%d"

func GetEntityFields(entityId *int64) (result structs.Fields, err error) {
	var pool *sql.DB = databasecache.Pool()
	rows, errQ := pool.Query(fmt.Sprintf(getEntityFields, *entityId))
	if errQ == nil {
		for rows.Next() {
			var field structs.Field
			rows.Scan(&field.ID, &field.ApiName, &field.Label, &field.Type, &field.Status, &field.Required, &field.ObjectId, &field.CreatedBy, &field.TextLength, &field.ObjectReference)
			result = append(result, field)
		}
		rows.Close()
	} else {
		err = errors.New("unable to run query against database to get entity fields")
		log.Print(errQ.Error(), errQ)
	}
	rows.Close()
	return
}

const alterAddColumnToTable = "ALTER TABLE `amper`.`%s` %s"
const insertFields = "INSERT INTO amper.field_sys VALUES %s"
const addForeginKey = "ALTER TABLE `amper`.`%s` ADD CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `amper`.`%s` (`identifier_sys`) ON DELETE NO ACTION ON UPDATE NO ACTION"

func CreateEntityField(userId *int64, entityId *int64, entityApiName *string, apiName *string, label *string, dataTypString *string, textLength *int64, objectReferenceId *int64, objectReferenceApiName *string, required *bool, status *bool) (bool, error) {
	dataType := constants.GetFieldType(dataTypString)
	if dataType == nil {
		return false, errors.New("data type is a required parameter and it has to reference a valid datatype")
	}
	var createStatement string
	if constants.GetDataTypes.REFERENCE.Name == *dataTypString {
		createStatement = fmt.Sprintf(constants.GetDataTypes.REFERENCE.CreateStatement, *apiName, "DEFAULT NULL")
		textLength = &constants.GetDataTypes.REFERENCE.Size
	} else if constants.GetDataTypes.TEXT.Name == *dataTypString {
		createStatement = fmt.Sprintf(constants.GetDataTypes.TEXT.CreateStatement, *apiName, *textLength, "DEFAULT NULL")
	} else {
		createStatement = fmt.Sprintf(dataType.CreateStatement, *apiName, "DEFAULT NULL")
	}

	var pool *sql.DB = databasecache.Pool()
	addColumnSQL := fmt.Sprintf(alterAddColumnToTable, *entityApiName, "ADD COLUMN "+createStatement)
	_, errDB := pool.Exec(addColumnSQL)
	if errDB != nil {
		log.Print(errDB.Error(), errDB)
		return false, errors.New("unable to create a database column with the provided parameters")
	}
	values := fmt.Sprintf("(null,  '%s', '%s', '%s', '%s', '%s', '%d', '%d', %s, %s)",
		*apiName, *label, *dataTypString,
		util.IfElse(status == nil || !*status, "0", "1").(string),
		util.IfElse(required == nil || !*required, "0", "1").(string),
		*entityId, *userId,
		util.IfElse(textLength != nil, util.IntToStr(textLength), "NULL").(string),
		util.IfElse(objectReferenceId != nil, util.IntToStr(objectReferenceId), "NULL").(string))

	insertFieldSQL := fmt.Sprintf(insertFields, values)
	_, errDB1 := pool.Exec(insertFieldSQL)
	if errDB1 != nil {
		log.Print(errDB.Error(), errDB)
		return false, errors.New("unable to create a database column with the provided parameters")
	}

	if constants.GetDataTypes.REFERENCE.Name == *dataTypString {
		addForeginKeyQuery := fmt.Sprintf(addForeginKey, *entityApiName, *apiName, *apiName, *objectReferenceApiName)
		_, errDB2 := pool.Exec(addForeginKeyQuery)
		if errDB2 != nil {
			log.Print(errDB2.Error(), errDB2)
			return false, errors.New("unable to create a database foregin key constraints with the provided parameters")
		}
	}
	return true, nil
}

const insertObjectTypeField = "INSERT INTO amper.object_type_field_sys VALUES (null, '%d', '%d', '%d', '%s')"

func AddObjectTypeField(userID *int64, fieldId *int64, objectTypeId *int64) (bool, error) {
	objectTypeField, errO := GetObjectTypeField(fieldId, &[]int64{*objectTypeId})
	if errO != nil {
		log.Print(errO.Error(), errO)
		return false, fmt.Errorf("unable to check of existing object type field for field: %d to object type: %d exist or not", *fieldId, *objectTypeId)
	} else if objectTypeField != nil {
		return false, fmt.Errorf("unable to add field: %d to object type: %d, the field already exist on the specified object type", *fieldId, *objectTypeId)
	}
	var pool *sql.DB = databasecache.Pool()
	timeFormatted := time.Now().Format("2006-01-02 15:04:05")
	_, errDb := pool.Exec(fmt.Sprintf(insertObjectTypeField, *objectTypeId, *fieldId, *userID, timeFormatted))
	if errDb != nil {
		log.Print(errDb.Error(), errDb)
		return false, fmt.Errorf("unable to add field: %d to object type: %d", *fieldId, *objectTypeId)
	}
	return true, nil
}

const deleteObjectTypeField = "DELETE from object_type_field_sys where object_type_id=%d and field_id=%d"

func DeleteObjectTypeField(userID *int64, fieldId *int64, objectTypeId *int64) (bool, error) {
	objectTypeField, errO := GetObjectTypeField(fieldId, &[]int64{*objectTypeId})
	if errO != nil {
		log.Print(errO.Error(), errO)
		return false, fmt.Errorf("unable to check of existing object type field for field: %d to object type: %d exist or not", *fieldId, *objectTypeId)
	} else if objectTypeField == nil {
		return false, fmt.Errorf("unable to delete field: %d to object type: %d, the field does not exist on the specified object type", *fieldId, *objectTypeId)
	}
	var pool *sql.DB = databasecache.Pool()
	res, errDb := pool.Exec(fmt.Sprintf(deleteObjectTypeField, *objectTypeId, *fieldId))
	if errDb != nil {
		log.Print(errDb.Error(), errDb)
		return false, fmt.Errorf("unable to delete field: %d to object type: %d, thre was a database communication error", *fieldId, *objectTypeId)
	} else if rowAf, _ := res.RowsAffected(); rowAf < 1 {
		return false, fmt.Errorf("unable to delete field: %d to object type: %d, no database rows were affected", *fieldId, *objectTypeId)
	}
	return true, nil
}

const getObjectTypeField = "SELECT id, object_type_id, field_id, created_date from amper.object_type_field_sys where object_type_id in (%s) and field_id=%d"

func GetObjectTypeField(fieldId *int64, objectTypeIds *[]int64) (result *structs.ObjectTypeField, err error) {
	var pool *sql.DB = databasecache.Pool()
	objectTypeIdsComma := strings.Join(arrays.IntToString(objectTypeIds), ",")
	rows, err := pool.Query(fmt.Sprintf(getObjectTypeField, objectTypeIdsComma, *fieldId))
	if err == nil {
		for rows.Next() {
			result = &structs.ObjectTypeField{}
			rows.Scan(&result.ID, &result.ObjectTypeId, &result.FieldId, &result.CreatedDate)
			rows.Close()
			break
		}
	} else {
		log.Print(err.Error(), err)
		return nil, fmt.Errorf("unable to lookup existing object type field for field: %d to object type: %s", *fieldId, objectTypeIdsComma)
	}
	return result, nil
}

const getObjectTypeFieldById = "SELECT id, object_type_id, field_id, created_date from amper.object_type_field_sys where field_id=%d"

func GetObjectTypeFieldById(fieldId *int64) (result *structs.ObjectTypeField, err error) {
	var pool *sql.DB = databasecache.Pool()
	rows, err := pool.Query(fmt.Sprintf(getObjectTypeFieldById, *fieldId))
	if err == nil {
		for rows.Next() {
			result = &structs.ObjectTypeField{}
			rows.Scan(&result.ID, &result.ObjectTypeId, &result.FieldId, &result.CreatedDate)
			rows.Close()
			break
		}
	} else {
		log.Print(err.Error(), err)
		return nil, fmt.Errorf("unable to lookup existing object type field for field: %d", *fieldId)
	}
	return result, nil
}
