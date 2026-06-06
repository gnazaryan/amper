package business

import (
	"amper/cache/business"
	"amper/common/structs"
	"amper/common/util"
	"amper/data/database"
	"fmt"
)

func GetUserRelationships(userId int64, employeeId int64) (result []structs.UserRelationshipExtended, err error) {
	result = make([]structs.UserRelationshipExtended, 0)
	userRelationships, errUR := database.GetUserRelationships(userId, employeeId)
	if errUR == nil {
		for _, userRelationship := range userRelationships {
			employee := business.GetUser(userRelationship.EmployeeId, true)
			manager := business.GetUser(userRelationship.ManagerId, true)
			result = append(result, structs.UserRelationshipExtended{
				ID:                userRelationship.ID,
				EmployeeId:        userRelationship.EmployeeId,
				EmployeeFirstName: employee.FirstName,
				EmployeeLastName:  employee.LastName,
				EmployeePhoto:     employee.Photo,
				ManagerId:         userRelationship.ManagerId,
				ManagerFirstName:  manager.FirstName,
				ManagerLastName:   manager.LastName,
				ManagerPhoto:      manager.Photo,
			})
		}
	} else {
		util.Loggify(errUR)
		err = fmt.Errorf("not able to fetch the user relationship, please try again later or contuct the support")
	}
	return result, err
}

func CreateUserRelationship(userId int64, employeeId int64, managerId int64) (result bool, err error) {
	_, errCUR := database.CreateUserRelationship(userId, employeeId, managerId)
	if errCUR != nil {
		util.Loggify(errCUR)
		err = fmt.Errorf("not able to create a user relationship, try again later or contuct the support")
		result = false
	} else {
		result = true
		//update the employees profile detail and add a manager
		employeeUser := business.GetUser(&employeeId, true)
		if employeeUser != nil && employeeUser.AmperId != nil {
			instance := business.GetAmperInstance(*employeeUser.AmperId)
			sessionId := GenerateSessionId(instance.Identifier, instance.Key, util.PointerString("app"))
			_, _, errPS := DedicatedCallWithRetry(&userId, &sessionId, map[string]string{
				"amperInstance": "profile/addRelationship",
			}, map[string]interface{}{
				"type":       "manager",
				"value":      managerId,
				"employeeId": employeeId,
			}, instance)
			if errPS != nil {
				result = false
			}
		} else {
			result = false
		}
		//update the managers profiledetail and add a reporter
		managerUser := business.GetUser(&managerId, true)
		if managerUser != nil && managerUser.AmperId != nil {
			instance := business.GetAmperInstance(*employeeUser.AmperId)
			sessionId := GenerateSessionId(instance.Identifier, instance.Key, util.PointerString("app"))
			_, _, errPS := DedicatedCallWithRetry(&userId, &sessionId, map[string]string{
				"amperInstance": "profile/addRelationship",
			}, map[string]interface{}{
				"type":       "reporter",
				"value":      employeeId,
				"employeeId": managerId,
			}, instance)
			if errPS != nil {
				result = false
			}
		} else {
			result = false
		}
	}
	return result, err
}

func DeleteUserRelationship(userId int64, managerId int64, employeeId int64) (result bool, err error) {
	success, errDRU := database.DeleteUserRelationship(userId, managerId, employeeId)
	if !success || errDRU != nil {
		util.Loggify(errDRU)
		err = fmt.Errorf("not able to delete a user relationship, try again later or contuct the support")
		result = false
	} else {
		result = true
		//update the employees profile detail and add a manager
		employeeUser := business.GetUser(&employeeId, true)
		if employeeUser != nil && employeeUser.AmperId != nil {
			instance := business.GetAmperInstance(*employeeUser.AmperId)
			sessionId := GenerateSessionId(instance.Identifier, instance.Key, util.PointerString("app"))
			_, _, errPS := DedicatedCallWithRetry(&userId, &sessionId, map[string]string{
				"amperInstance": "profile/removeRelationship",
			}, map[string]interface{}{
				"type":       "manager",
				"value":      managerId,
				"employeeId": employeeId,
			}, instance)
			if errPS != nil {
				result = false
			}
		} else {
			result = false
		}
		//update the managers profiledetail and add a reporter
		managerUser := business.GetUser(&managerId, true)
		if managerUser != nil && managerUser.AmperId != nil {
			instance := business.GetAmperInstance(*employeeUser.AmperId)
			sessionId := GenerateSessionId(instance.Identifier, instance.Key, util.PointerString("app"))
			_, _, errPS := DedicatedCallWithRetry(&userId, &sessionId, map[string]string{
				"amperInstance": "profile/removeRelationship",
			}, map[string]interface{}{
				"type":       "reporter",
				"value":      employeeId,
				"employeeId": managerId,
			}, instance)
			if errPS != nil {
				result = false
			}
		} else {
			result = false
		}
	}
	return result, err
}
