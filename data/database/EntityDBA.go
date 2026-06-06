package database

import (
	databasecache "amper/cache/database"
	"amper/common/constants"
	"amper/common/structs"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

const getEntitiesWithPagination = "SELECT * FROM object_sys LIMIT %d OFFSET %d"

// GetEntities query the database and retrieve the entities for the provided start and limit parameters
func GetEntities(start *int, limit *int) (result []structs.Entity, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(getEntitiesWithPagination, *limit, *start)
	rows, errQ := pool.Query(query)
	if errQ == nil {
		for rows.Next() {
			var entity structs.Entity
			rows.Scan(&entity.ID, &entity.ApiName, &entity.Title, &entity.TitlePlural)
			result = append(result, entity)
		}
	} else {
		err = errors.New("unable to run query against database to get entities")
		log.Print(errQ.Error(), errQ)
	}
	rows.Close()
	return
}

const GET_ENTITY_BY_ID = "SELECT * FROM amper.object_sys where id='%d'"

func GetEntity(EntityId *int64) (result *structs.Entity, err error) {
	var pool *sql.DB = databasecache.Pool()
	rows, errQ := pool.Query(fmt.Sprintf(GET_ENTITY_BY_ID, *EntityId))
	if errQ == nil {
		for rows.Next() {
			result = new(structs.Entity)
			rows.Scan(&result.ID, &result.ApiName, &result.Title, &result.TitlePlural)
			rows.Close()
			return
		}
	} else {
		err = errors.New("unable to run query against database to get entity")
		log.Print(errQ.Error(), errQ)
	}
	return
}

const GET_ENTITY_BY_API_NAME = "SELECT * FROM amper.object_sys where apiName='%s'"

func GetEntityByApiName(ApiName *string) (result *structs.Entity, err error) {
	var pool *sql.DB = databasecache.Pool()
	rows, errQ := pool.Query(fmt.Sprintf(GET_ENTITY_BY_API_NAME, *ApiName))
	if errQ == nil {
		for rows.Next() {
			result = new(structs.Entity)
			rows.Scan(&result.ID, &result.ApiName, &result.Title, &result.TitlePlural)
			rows.Close()
			return
		}
	} else {
		err = errors.New("unable to run query against database to get entity")
		log.Print(errQ.Error(), errQ)
		rows.Close()
	}
	return
}

func DeleteEntity(EntityId *int64) (result bool, err error) {
	entity, errDB := GetEntity(EntityId)
	if errDB == nil && entity != nil {
		var pool *sql.DB = databasecache.Pool()

		if deleteObjectTypeFields(pool, entity.ID) &&
			deleteObjectType(pool, entity.ID) &&
			removeFieldsWithEntityId(pool, entity.ID) &&
			removeObjectWithApiName(pool, entity.ApiName) &&
			dropTableWithApiName(pool, entity.ApiName) {
			result = true
		} else {
			result = false
			err = fmt.Errorf("unable to delete an object entry with specified object id: %d", *EntityId)
		}
	} else {
		result = false
		err = fmt.Errorf("unable to locate an object entry with specified object id: %d", *EntityId)
	}
	return
}

const INSERT_TABLE = "INSERT INTO amper.object_sys VALUES (null, '%s', '%s', '%s')"
const GET_TABLE = "SELECT id FROM amper.object_sys where apiName='%s'"
const GET_OBJECT_TYPE = "SELECT id FROM amper.object_type_sys where object_id=%d and apiName='%s'"
const GET_FIELDS = "SELECT id FROM amper.field_sys where entityId=%d"
const INSERT_OBJECT_TYPE_FIELD = "INSERT INTO amper.object_type_field_sys VALUES (null, '%d', '%d', '%d', '%s')"

func CreateEntity(userId *int64, ApiName *string, Title *string, TitlePlural *string) (success bool, err error) {
	//Create Database Table with specified api name and default fields
	createTableStatement := buildCreateTable(ApiName)
	var pool *sql.DB = databasecache.Pool()
	_, errDb := pool.Exec(createTableStatement)
	if errDb == nil {
		insertTableStatement := fmt.Sprintf(INSERT_TABLE, *ApiName, *Title, *TitlePlural)
		_, errDb = pool.Exec(insertTableStatement)
		if errDb == nil {
			var id int64
			errDb := pool.QueryRow(fmt.Sprintf(GET_TABLE, *ApiName)).Scan(&id)
			if errDb == nil {
				insertFieldStatement := buildInsertFields(&id)
				_, errDb = pool.Exec(insertFieldStatement)
				if errDb == nil {
					rows, errDb := pool.Query(fmt.Sprintf(GET_FIELDS, id))
					if errDb == nil {
						fieldIds := []int{}
						for rows.Next() {
							var fieldId int
							rows.Scan(&fieldId)
							fieldIds = append(fieldIds, fieldId)
						}
						insertObjectTypeStatement := buildInsertObjectType(userId, &id)
						_, errDb = pool.Exec(insertObjectTypeStatement)
						if errDb == nil {
							var objectTypeId int64
							errDb := pool.QueryRow(fmt.Sprintf(GET_OBJECT_TYPE, id, constants.GetBaseObjectType().KEY)).Scan(&objectTypeId)
							if errDb == nil {
								for _, fieldId := range fieldIds {
									timeFormatted := time.Now().Format("2006-01-02 15:04:05")
									objectTypeFieldStatement := fmt.Sprintf(INSERT_OBJECT_TYPE_FIELD, objectTypeId, fieldId, *userId, timeFormatted)
									_, errDb := pool.Exec(objectTypeFieldStatement)
									if errDb == nil {
										success = true
									} else {
										deleteObjectTypeFieldsByObjectTypeId(pool, &objectTypeId)
										deleteObjectType(pool, &id)
										removeFieldsWithEntityId(pool, &id)
										removeObjectWithApiName(pool, ApiName)
										dropTableWithApiName(pool, ApiName)
										success = false
										err = fmt.Errorf("unable to insert a object type field entry with specified parameters and apiname: %s", *ApiName)
										return
									}
								}
							} else {
								deleteObjectType(pool, &id)
								removeFieldsWithEntityId(pool, &id)
								removeObjectWithApiName(pool, ApiName)
								dropTableWithApiName(pool, ApiName)
								success = false
								err = fmt.Errorf("unable to insert a object type entry with specified parameters and apiname: %s", *ApiName)
							}
						} else {
							removeFieldsWithEntityId(pool, &id)
							removeObjectWithApiName(pool, ApiName)
							dropTableWithApiName(pool, ApiName)
							success = false
							err = fmt.Errorf("unable to insert a object type entry with specified parameters and apiname: %s", *ApiName)
						}
					} else {
						removeFieldsWithEntityId(pool, &id)
						removeObjectWithApiName(pool, ApiName)
						dropTableWithApiName(pool, ApiName)
						success = false
						err = fmt.Errorf("unable to insert a object fields entries with specified parameters and apiname: %s", *ApiName)
					}
					rows.Close()
				} else {
					removeFieldsWithEntityId(pool, &id)
					removeObjectWithApiName(pool, ApiName)
					dropTableWithApiName(pool, ApiName)
					success = false
					err = fmt.Errorf("unable to insert a object fields entries with specified parameters and apiname: %s", *ApiName)
				}
			} else {
				removeObjectWithApiName(pool, ApiName)
				dropTableWithApiName(pool, ApiName)
				success = false
				err = fmt.Errorf("unable to insert a table entry with specified parameters and apiname: %s", *ApiName)
			}
		} else {
			removeObjectWithApiName(pool, ApiName)
			dropTableWithApiName(pool, ApiName)
			success = false
			err = fmt.Errorf("unable to insert a table entry with specified parameters and apiname: %s", *ApiName)
		}
	} else {
		dropTableWithApiName(pool, ApiName)
		success = false
		err = fmt.Errorf("unable to create a table with specified parameters and apiname: %s", *ApiName)
	}
	return
}

const DELETE_OBJECT_TYPE_FIELDS = "DELETE OTF from amper.object_type_field_sys as OTF INNER JOIN amper.object_type_sys as OT ON OT.id=OTF.object_type_id AND OT.object_id='%d'"

func deleteObjectTypeFields(pool *sql.DB, entityId *int64) (result bool) {
	if entityId != nil {
		_, errDB := pool.Exec(fmt.Sprintf(DELETE_OBJECT_TYPE_FIELDS, *entityId))
		if errDB != nil {
			log.Print(errDB.Error(), errDB)
		}
		return errDB == nil
	}
	return true
}

const DELETE_OBJECT_TYPE_FIELDS_BY_OBJECT_TYPE_ID = "DELETE from amper.object_type_field_sys WHERE object_type_id='%d'"

func deleteObjectTypeFieldsByObjectTypeId(pool *sql.DB, objectTypeId *int64) (result bool) {
	if objectTypeId != nil {
		_, errDB := pool.Exec(fmt.Sprintf(DELETE_OBJECT_TYPE_FIELDS_BY_OBJECT_TYPE_ID, *objectTypeId))
		if errDB != nil {
			log.Print(errDB.Error(), errDB)
		}
		return errDB == nil
	}
	return true
}

const DELETE_OBJECT_TYPE = "DELETE from amper.object_type_sys where object_id='%d'"

func deleteObjectType(pool *sql.DB, entityId *int64) (result bool) {
	_, errDB := pool.Exec(fmt.Sprintf(DELETE_OBJECT_TYPE, *entityId))
	if errDB != nil {
		log.Print(errDB.Error(), errDB)
	}
	return errDB == nil
}

const DROP_TABLE = "DROP TABLE `amper`.`%s`"

func dropTableWithApiName(pool *sql.DB, ApiName *string) (result bool) {
	_, errDB := pool.Exec(fmt.Sprintf(DROP_TABLE, *ApiName))
	if errDB != nil {
		log.Print(errDB.Error(), errDB)
	}
	return errDB == nil
}

const REMOVE_ENTITY = "DELETE from amper.object_sys where apiName='%s'"

func removeObjectWithApiName(pool *sql.DB, ApiName *string) (result bool) {
	_, errDB := pool.Exec(fmt.Sprintf(REMOVE_ENTITY, *ApiName))
	if errDB != nil {
		log.Print(errDB.Error(), errDB)
	}
	return errDB == nil
}

const DELETE_OBJECT_FIELDS = "DELETE from amper.field_sys where entityId=%d"

func removeFieldsWithEntityId(pool *sql.DB, entityId *int64) (result bool) {
	_, errDB := pool.Exec(fmt.Sprintf(DELETE_OBJECT_FIELDS, *entityId))
	if errDB != nil {
		log.Print(errDB.Error(), errDB)
	}
	return errDB == nil
}

const INSERT_OBJECT_TYPE = "INSERT INTO amper.object_type_sys VALUES (null, '%d', '%s', '%s', %s, '%d', '%s')"

func buildInsertObjectType(userId *int64, entityId *int64) (result string) {
	timeFormatted := time.Now().Format("2006-01-02 15:04:05")
	result = fmt.Sprintf(INSERT_OBJECT_TYPE, *entityId, constants.GetBaseObjectType().KEY, constants.GetBaseObjectType().LABEL, "null", *userId, timeFormatted)
	return
}

const CREATE_TABLE = "CREATE TABLE `amper`.`%s` (%s, PRIMARY KEY (`id`), UNIQUE KEY `identifier_sys_UNIQUE` (`identifier_sys`)) AUTO_INCREMENT=1"

func buildCreateTable(ApiName *string) (result string) {
	createStatements := [...]string{
		fmt.Sprintf(constants.GetDefaultFields().ID.Type.CreateStatement, constants.GetDefaultFields().ID.ApiName),
		fmt.Sprintf(constants.GetDefaultFields().IDENTIFIER.Type.CreateStatement, constants.GetDefaultFields().IDENTIFIER.ApiName, constants.GetDefaultFields().IDENTIFIER.Type.Size, "NOT NULL"),
		fmt.Sprintf(constants.GetDefaultFields().OBJECTTYPE.Type.CreateStatement, constants.GetDefaultFields().OBJECTTYPE.ApiName, constants.GetDefaultFields().OBJECTTYPE.Type.Size, "NOT NULL"),
		fmt.Sprintf(constants.GetDefaultFields().NAME.Type.CreateStatement, constants.GetDefaultFields().NAME.ApiName, constants.GetDefaultFields().NAME.Type.Size, "NOT NULL"),
		fmt.Sprintf(constants.GetDefaultFields().STATUS.Type.CreateStatement, constants.GetDefaultFields().STATUS.ApiName, "NOT NULL"),
	}

	result = fmt.Sprintf(CREATE_TABLE, *ApiName, strings.Join(createStatements[:], ","))
	return
}

const INSERT_FIELDS = "INSERT INTO amper.field_sys VALUES %s"

func buildInsertFields(entityId *int64) (result string) {
	insertStatements := [...]string{
		fmt.Sprintf("(null,  '%s', '%s', '%s', '%d', '%d', '%d', '1', NULL, NULL)", constants.GetDefaultFields().ID.ApiName,
			constants.GetDefaultFields().ID.Label, constants.GetDefaultFields().ID.Type.Name, constants.GetDefaultFields().ID.Status, constants.GetDefaultFields().ID.Required, *entityId),
		fmt.Sprintf("(null,  '%s', '%s', '%s', '%d', '%d', '%d', '1', '%d', NULL)", constants.GetDefaultFields().IDENTIFIER.ApiName,
			constants.GetDefaultFields().IDENTIFIER.Label, constants.GetDefaultFields().IDENTIFIER.Type.Name, constants.GetDefaultFields().IDENTIFIER.Status, constants.GetDefaultFields().IDENTIFIER.Required, *entityId, constants.GetDefaultFields().IDENTIFIER.Type.Size),
		fmt.Sprintf("(null,  '%s', '%s', '%s', '%d', '%d', '%d', '1', '%d', NULL)", constants.GetDefaultFields().OBJECTTYPE.ApiName,
			constants.GetDefaultFields().OBJECTTYPE.Label, constants.GetDefaultFields().OBJECTTYPE.Type.Name, constants.GetDefaultFields().OBJECTTYPE.Status, constants.GetDefaultFields().OBJECTTYPE.Required, *entityId, constants.GetDefaultFields().OBJECTTYPE.Type.Size),
		fmt.Sprintf("(null,  '%s', '%s', '%s', '%d', '%d', '%d', '1', '%d', NULL)", constants.GetDefaultFields().NAME.ApiName,
			constants.GetDefaultFields().NAME.Label, constants.GetDefaultFields().NAME.Type.Name, constants.GetDefaultFields().NAME.Status, constants.GetDefaultFields().NAME.Required, *entityId, constants.GetDefaultFields().NAME.Type.Size),
		fmt.Sprintf("(null,  '%s', '%s', '%s', '%d', '%d', '%d', '1', NULL, NULL)", constants.GetDefaultFields().STATUS.ApiName,
			constants.GetDefaultFields().STATUS.Label, constants.GetDefaultFields().STATUS.Type.Name, constants.GetDefaultFields().STATUS.Status, constants.GetDefaultFields().STATUS.Required, *entityId),
	}
	result = fmt.Sprintf(INSERT_FIELDS, strings.Join(insertStatements[:], ","))
	return
}

const UPDATE_TABLE = "UPDATE amper.object_sys SET title='%s', titlePlural='%s' WHERE id=%d and apiName='%s'"

func EditEntity(EntityId *int64, ApiName *string, Title *string, TitlePlural *string) (result bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	resultDB, errDB := pool.Exec(fmt.Sprintf(UPDATE_TABLE, *Title, *TitlePlural, *EntityId, *ApiName))
	if rows, errRow := resultDB.RowsAffected(); errDB == nil && errRow == nil && rows > 0 {
		result = true
	} else {
		result = false
		err = fmt.Errorf("unable to modify an object with specified parameters and apiname: %s", *ApiName)
	}
	return
}

const GET_FIELDS_BY_ENTITY_ID = "SELECT * FROM amper.field_sys where entityId=%d"

func GetFields(UserId *int64, ObjectId *int64) (result *[]structs.Field, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(GET_FIELDS_BY_ENTITY_ID, *ObjectId)
	rows, errQ := pool.Query(query)
	if errQ == nil {
		result = new([]structs.Field)
		for rows.Next() {
			var field structs.Field
			rows.Scan(&field.ID, &field.ApiName, &field.Label, &field.Type, &field.Status, &field.Required, &field.ObjectId, &field.CreatedBy, &field.TextLength, &field.ObjectReference)
			*result = append(*result, field)
		}
	} else {
		err = errors.New("unable to run query against database to get fields")
		log.Print(errQ.Error(), errQ)
	}
	rows.Close()
	return
}
