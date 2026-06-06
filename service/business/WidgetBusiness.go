package business

import (
	"amper/common/constants"
	"amper/common/structs"
	"amper/common/util/jsons"
	"amper/data/database"
	"fmt"
	"log"
	"reflect"
)

func GetWidgets(userId *int64, dashboardId *int64) (result *structs.Dashboard, err error) {
	result, errDB := database.GetWidgets(userId, dashboardId)
	if errDB != nil {
		log.Print(errDB.Error(), errDB)
		err = fmt.Errorf("unable to fetch dashboard widgets, please try again later or contact the support")
	}
	return
}

func GetInteractions(userId *int64, dashboardId *int64, widgetId *int64, objectApiName *string) (result *structs.Dashboard, err error) {
	result, errDB := database.GetWidgets(userId, dashboardId)
	if errDB != nil {
		log.Print(errDB.Error(), errDB)
		return nil, fmt.Errorf("unable to fetch dashboard widget interactions, please try again later or contact the support")
	}
	widgets := make([]structs.Widget, 0)
	var currentWidget *structs.Widget
	for _, widget := range *result.Widgets {
		if *widgetId == *widget.ID {
			currentWidget = &widget
			break
		}
	}
	if currentWidget != nil {
		configurationJson := (*currentWidget).Configuration
		if configurationJson != nil {
			configurationJson := currentWidget.Configuration
			if configurationJson != nil {
				configuration, errJ := jsons.GetJsonObject(configurationJson)
				if errJ == nil && reflect.TypeOf(configuration).String() == "map[string]interface {}" {
					widgetType, hasType := configuration["type"].(string)
					if hasType && widgetType == "recordList" {
						if objectApiName == nil || len(*objectApiName) < 1 {
							return nil, fmt.Errorf("api name is required parameter for calculating the widget interactions")
						}
						object, errO := GetEntityByApiName(userId, objectApiName)
						if errO != nil {
							return nil, fmt.Errorf("no object found with supplied object api name: %s", *objectApiName)
						}
						fields, errF := GetFields(userId, object.ID)
						if errF != nil {
							return nil, fmt.Errorf("no object fields found with supplied object api name: %s", *objectApiName)
						}
						for _, widget := range *result.Widgets {
							if *widgetId != *widget.ID {
								configurationJson := widget.Configuration
								if configurationJson != nil {
									configuration, errJ := jsons.GetJsonObject(configurationJson)
									if errJ == nil && reflect.TypeOf(configuration).String() == "map[string]interface {}" {
										object, hasIOW := configuration["object"].(map[string]interface{})
										if hasIOW && object != nil && object["apiName"] != nil {
											apiName := object["apiName"].(string)
											if len(apiName) > 0 {
												widgetObject, errOW := GetEntityByApiName(userId, &apiName)
												if errOW == nil {
													for _, field := range *fields {
														if *field.Type == constants.GetDataTypes.REFERENCE.Name &&
															*field.ObjectReference == *widgetObject.ID {
															widgets = append(widgets, widget)
														}
													}
												}
											}
										}
									}
								}
							}
						}
					} else if hasType && widgetType == "recordDetail" {
						for _, widget := range *result.Widgets {
							if *widgetId != *widget.ID {
								configurationJson := widget.Configuration
								if configurationJson != nil {
									configuration, errJ := jsons.GetJsonObject(configurationJson)
									if errJ == nil && reflect.TypeOf(configuration).String() == "map[string]interface {}" {
										widgetType, hasType := configuration["type"].(string)
										if hasType && widgetType == "recordList" {
											widgets = append(widgets, widget)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	result.Widgets = &widgets
	return result, nil
}

func AddWidget(userId *int64, dashboardId *int64, label *string, description *string, configuration *string) (result bool, err error) {
	result, errDB := database.AddWidget(userId, dashboardId, label, description, configuration)
	if errDB != nil {
		log.Print(errDB.Error(), errDB)
		err = fmt.Errorf("unable to add a widget to dashboard with specified parameters, please try again later or contact the support")
	}
	return
}

func RemoveWidget(userId *int64, dashboardId *int64, widgetId *int64) (result bool, err error) {
	result, errDB := database.RemoveWidget(userId, dashboardId, widgetId)
	if errDB != nil {
		log.Print(errDB.Error(), errDB)
		err = fmt.Errorf("unable to remove a dashboard widget, please try again later or contact the support")
	}
	return
}

func UpdateWidget(userId *int64, dashboard *structs.Dashboard) (bool, error) {
	for _, widget := range *dashboard.Widgets {
		widgetRes, errW := database.UpdateWidget(userId, dashboard.ID, widget.ID, widget.Label, widget.Description, widget.Configuration)
		if !widgetRes || errW != nil {
			if errW != nil {
				log.Print(errW.Error(), errW)
			}
			return false, fmt.Errorf("unable to update a widget to dashboard with specified parameters, please try again later or contact the support")
		}
	}
	return true, nil
}
