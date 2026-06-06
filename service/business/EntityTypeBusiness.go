package business

import (
	"amper/common/constants"
	"amper/common/structs"
	"amper/data/database"
	"fmt"
	"log"
	"sort"
	"time"
)

func GetEntityTypesMap(userId *int64, objectId *int64) (map[int64]structs.ObjectType, error) {
	result := make(map[int64]structs.ObjectType)
	objectTypes, err := GetObjectTypes(userId, objectId)
	if err == nil && len(*objectTypes) > 0 {
		for _, objType := range *objectTypes {
			result[*objType.ID] = objType
		}
	}
	return result, err
}

// GetEntities is responsible for retrieving entities with the provided start and limit parameters
func GetObjectTypes(userID *int64, objectId *int64) (result *[]structs.ObjectType, err error) {
	result, err = database.GetObjectTypesFields(userID, objectId)
	if err != nil {
		err = fmt.Errorf("unable to retrieve object types for %d", *objectId)
		log.Print(err.Error(), err)
	}
	//Sort the result based on the object type create date
	sort.Slice(*result, func(i, j int) bool {
		first, _ := time.Parse(constants.TIME_FORMAT, *(*result)[i].CreatedDate)
		second, _ := time.Parse(constants.TIME_FORMAT, *(*result)[j].CreatedDate)
		return first.Before(second)
	})

	return
}

// GetEntities is responsible for retrieving entities with the provided start and limit parameters
func GetObjectTypeFields(userID *int64, objectId *int64, ObjectTyoeId *int64) (result *structs.ObjectType, err error) {
	result, err = database.GetObjectTypeFields(userID, objectId, ObjectTyoeId)
	if err != nil {
		err = fmt.Errorf("unable to retrieve object types for %d", *objectId)
		log.Print(err.Error(), err)
	}
	return
}

func CreateObjectType(UserID *int64, ObjectId *int64, ExtendsTo *int64, ApiName *string, Label *string) (result bool, err error) {
	entity, err := database.GetEntity(ObjectId)
	if err == nil && entity.ID != nil {
		objectType, errOT := database.GetObjectTypeByApiName(ObjectId, ApiName)
		if objectType == nil && errOT == nil {
			result, err = database.CreateObjectType(UserID, ObjectId, ExtendsTo, ApiName, Label)
			if err != nil {
				log.Print(err.Error(), err)
				err = fmt.Errorf("unable to create object type for api name: %s", *ApiName)
			}
		} else if objectType != nil {
			err = fmt.Errorf("unable to create object type for api name: %s, object with specified name already exists", *ApiName)
			log.Print(err.Error(), err)
		} else {
			err = fmt.Errorf("unable to identify if object with specified Api Name: %s already exists", *ApiName)
			log.Print(errOT.Error(), errOT)
		}
	} else {
		err = fmt.Errorf("unable to create object type for the provided ApiName: %s, there is no object found with the provided id %d", *ApiName, *ObjectId)
		log.Print(err.Error(), err)
	}
	return
}

func DeleteObjectType(UserID *int64, ObjectTypeId *int64) (result bool, err error) {
	result, err = database.DeleteObjectType(UserID, ObjectTypeId)
	if err != nil {
		err = fmt.Errorf("unable to delete object type for id: %d", *ObjectTypeId)
		log.Print(err.Error(), err)
	}
	return
}
