package business

import (
	"amper/common/structs"
	"amper/data/database"
	"fmt"
	"log"
)

// GetEntities is responsible for retrieving entities with the provided start and limit parameters
func GetEntities(start *int, limit *int) (entities []structs.Entity, err error) {
	entities, errDb := database.GetEntities(start, limit)
	if errDb != nil {
		err = fmt.Errorf("unable to retrieve entities for the provided start: %d and limit: %d", start, limit)
		log.Print(errDb.Error(), errDb)
	}
	return
}

// GetEntity is for geting an existing entity with the given id
func GetEntity(userID *int64, EntityId *int64) (entity *structs.Entity, err error) {
	entity, errDb := database.GetEntity(EntityId)
	if errDb != nil {
		err = fmt.Errorf("unable to retrieve entity for the provided EntityId: %d", EntityId)
		log.Print(errDb.Error(), errDb)
	}
	return
}

// GetEntity is for geting an existing entity with the given id
func GetEntityByApiName(userID *int64, apiName *string) (entity *structs.Entity, err error) {
	entity, errDb := database.GetEntityByApiName(apiName)
	if errDb != nil {
		err = fmt.Errorf("unable to retrieve entity for the provided api name: %s", *apiName)
		log.Print(errDb.Error(), errDb)
	}
	return
}

// DeleteEntity is for deleting an existing entity with given id
func DeleteEntity(userID *int64, EntityId *int64) (result bool, err error) {
	result, errDb := database.DeleteEntity(EntityId)
	if errDb != nil {
		err = fmt.Errorf("unable to delete entity for the provided EntityId: %d", EntityId)
		log.Print(errDb.Error(), errDb)
	}
	return
}

// CreateEntity is for creting an entity with the provided information
func CreateEntity(userID *int64, ApiName *string, Title *string, TitlePlural *string) (result bool, err error) {
	entity, err := database.GetEntityByApiName(ApiName)
	if err == nil && entity == nil {
		result, err = database.CreateEntity(userID, ApiName, Title, TitlePlural)
		if err != nil {
			err = fmt.Errorf("unable to create entity for the provided ApiName: %s", *ApiName)
			log.Print(err.Error(), err)
		}
	} else {
		err = fmt.Errorf("unable to create entity for the provided ApiName: %s, there is already an object with the provided details", *ApiName)
		log.Print(err.Error(), err)
	}
	return
}

// EditEntity is for editing an entity with the provided id and information
func EditEntity(userID *int64, EntityId *int64, ApiName *string, Title *string, TitlePlural *string) (result bool, err error) {
	entity, err := database.GetEntityByApiName(ApiName)
	if err == nil && entity != nil {
		success, errDb := database.EditEntity(EntityId, ApiName, Title, TitlePlural)
		if errDb != nil && !success {
			err = fmt.Errorf("unable to edit entity for the provided EntityId: %d", EntityId)
			log.Print(errDb.Error(), errDb)
		} else {
			result = true
		}
	} else {
		err = fmt.Errorf("unable to create entity for the provided ApiName: %s, there is already an object with the provided details", *ApiName)
		log.Print(err.Error(), err)
	}
	return
}

func GetFields(UserId *int64, ObjectId *int64) (result *[]structs.Field, err error) {
	result, err = database.GetFields(UserId, ObjectId)
	if err != nil {
		err = fmt.Errorf("unable to retrieve fields for the provided object: %d, please try later", *ObjectId)
		log.Print(err.Error(), err)
	}
	return
}
