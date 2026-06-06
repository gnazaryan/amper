package authorization

import (
	"amper/common/argument"
	"amper/common/structs"
	"amper/service/business"
)

func GetWidgets(userId *int64, dashboardId *int64) (result *structs.Dashboard, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "dashboardId": dashboardId})
	if err != nil {
		return nil, err
	}
	result, err = business.GetWidgets(userId, dashboardId)
	return
}

func GetInteractions(userId *int64, dashboardId *int64, widgetId *int64, objectApiName *string) (result *structs.Dashboard, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "dashboardId": dashboardId, "widgetId": widgetId})
	if err != nil {
		return nil, err
	}
	result, err = business.GetInteractions(userId, dashboardId, widgetId, objectApiName)
	return
}

func AddWidget(userId *int64, dashboardId *int64, label *string, description *string, configuration *string) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "dashboardId": dashboardId, "label": label, "description": label, "configuration": label})
	if err != nil {
		return false, err
	}
	result, err = business.AddWidget(userId, dashboardId, label, description, configuration)
	return
}

func RemoveWidget(userId *int64, dashboardId *int64, widgetId *int64) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "dashboardId": dashboardId, "widgetId": widgetId})
	if err != nil {
		return false, err
	}
	result, err = business.RemoveWidget(userId, dashboardId, widgetId)
	return
}

func UpdateWidget(userId *int64, dashboard *structs.Dashboard) (result bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "dashboardId": dashboard.ID, "Widgets": dashboard.Widgets, "widgetId": (*dashboard.Widgets)[0].ID, "label": (*dashboard.Widgets)[0].Label, "description": (*dashboard.Widgets)[0].Description, "configuration": (*dashboard.Widgets)[0].Configuration})
	if err != nil {
		return false, err
	}
	result, err = business.UpdateWidget(userId, dashboard)
	return
}
