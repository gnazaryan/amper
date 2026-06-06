package authorization

import (
	"amper/common/argument"
	"amper/common/structs"
	"amper/service/business"
)

func GetDashboards(userId *int64) (result *[]structs.Dashboard, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId})
	if err != nil {
		return nil, err
	}
	result, err = business.GetDashboards(userId)
	return
}

func AddDashboard(userId *int64, label *string, description *string, configuration *string) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "label": label, "description": description, "configuration": configuration})
	if err != nil {
		return false, err
	}
	result, err = business.AddDashboard(userId, label, description, configuration)
	return
}

func UpdateDashboard(userId *int64, id *int64, label *string, description *string, configuration *string) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "id": id, "label": label, "description": description, "configuration": configuration})
	if err != nil {
		return false, err
	}
	result, err = business.UpdateDashboard(userId, id, label, description, configuration)
	return
}

func RemoveDashboard(userId *int64, id *int64) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userId, "dashboardId": id})
	if err != nil {
		return false, err
	}
	result, err = business.RemoveDashboard(userId, id)
	return
}
