package business

import (
	"amper/common/structs"
	"amper/data/database"
	"fmt"
	"log"
)

func GetDashboards(userId *int64) (result *[]structs.Dashboard, err error) {
	result, errDB := database.GetDashboards(userId)
	if errDB != nil {
		log.Print(errDB.Error(), errDB)
		err = fmt.Errorf("unable to fetch dashboards, please try again later or contact the support")
	}
	return
}

func AddDashboard(userId *int64, label *string, description *string, configuration *string) (result bool, err error) {
	result, errDB := database.AddDashboard(userId, label, description, configuration)
	if errDB != nil {
		log.Print(errDB.Error(), errDB)
		err = fmt.Errorf("unable to add a dashboard with specified parameters, please try again later or contact the support")
	}
	return
}

func UpdateDashboard(userId *int64, id *int64, label *string, description *string, configuration *string) (result bool, err error) {
	result, errDB := database.UpdateDashboard(userId, id, label, description, configuration)
	if errDB != nil {
		log.Print(errDB.Error(), errDB)
		err = fmt.Errorf("unable to update a dashboard with specified parameters, please try again later or contact the support")
	}
	return
}

func RemoveDashboard(userId *int64, id *int64) (result bool, err error) {
	result, errDB := database.RemoveDashboard(userId, id)
	if errDB != nil {
		log.Print(errDB.Error(), errDB)
		err = fmt.Errorf("unable to remove a dashboard, please try again later or contact the support")
	}
	return
}
