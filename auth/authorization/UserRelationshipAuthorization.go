package authorization

import (
	"amper/common/argument"
	"amper/common/structs"
	"amper/service/business"
)

func GetUserRelationships(userId int64, employeeId int64) (result []structs.UserRelationshipExtended, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "employeeId": employeeId})
	if err != nil {
		return nil, err
	}
	return business.GetUserRelationships(userId, employeeId)
}

func CreateUserRelationship(userId int64, employeeId int64, managerId int64) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "employeeId": employeeId, "managerId": managerId})
	if err != nil {
		return false, err
	}
	return business.CreateUserRelationship(userId, employeeId, managerId)
}

func DeleteUserRelationship(userId int64, managerId int64, employeeId int64) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "managerId": managerId, "employeeId": employeeId})
	if err != nil {
		return false, err
	}
	return business.DeleteUserRelationship(userId, managerId, employeeId)
}
