package database

import (
	databasecache "amper/cache/database"
	"amper/common/structs"
	"amper/common/util/datetime"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

const getDashboardWidgets = "SELECT * FROM amper.widget_sys WHERE created_by=%d AND dashboard_id=%d"

func GetWidgets(userId *int64, dashboardId *int64) (result *structs.Dashboard, err error) {
	dashboard, errD := GetDashboard(userId, dashboardId)
	if errD != nil {
		err = errors.New("unable to run query against database to get dashboard")
		log.Print(errD.Error(), errD)
	}
	var pool *sql.DB = databasecache.Pool()
	rows, errQ := pool.Query(fmt.Sprintf(getDashboardWidgets, *userId, *dashboardId))
	var widgets []structs.Widget
	if errQ == nil {
		for rows.Next() {
			var widget structs.Widget
			var dashboardId int64
			rows.Scan(&widget.ID, &dashboardId, &widget.Label, &widget.Description, &widget.Configuration, &widget.CreatedDate, &widget.CreatedBy)
			widgets = append(widgets, widget)
		}
	} else {
		err = errors.New("unable to run query against database to get dashboard widgets")
		log.Print(errQ.Error(), errQ)
	}
	rows.Close()
	dashboard.Widgets = &widgets
	result = &dashboard
	return
}

const insertDashboardWidget = "INSERT INTO amper.widget_sys VALUES (null, %d, '%s', '%s', '%s', '%s', %d)"

func AddWidget(userId *int64, dashboardId *int64, label *string, description *string, configuration *string) (result bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(insertDashboardWidget, *dashboardId, *label, *description, *configuration, datetime.GetDateTimeFormatted(), *userId)
	res, errW := pool.Exec(query)
	if errW != nil {
		log.Print(errW.Error(), errW)
		err = fmt.Errorf("unable to inser a dashboard widget with specified parameters: label - %s, description - %s, configuration - %s", *label, *description, *configuration)
		return
	} else if count, _ := res.RowsAffected(); count < 1 {
		err = fmt.Errorf("unable to inser a dashboard widget with specified parameters: label - %s, description - %s, configuration - %s", *label, *description, *configuration)
	}
	result = true
	return
}

const deleteDashboardWidget = "DELETE from amper.widget_sys where dashboard_id=%d AND id=%d AND created_by=%d"

func RemoveWidget(userId *int64, dashboardId *int64, widgetId *int64) (result bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(deleteDashboardWidget, *dashboardId, *widgetId, *userId)
	res, errW := pool.Exec(query)
	if errW != nil {
		log.Print(errW.Error(), errW)
		err = fmt.Errorf("unable to delete a dashboard with specified id - %d", *dashboardId)
		return
	} else if count, _ := res.RowsAffected(); count < 1 {
		err = fmt.Errorf("unable to delete a dashboard with specified id - %d", *dashboardId)
	}
	result = true
	return
}

const deleteDashboardWidgets = "DELETE from amper.widget_sys where dashboard_id=%d AND created_by=%d"

func RemoveWidgets(userId *int64, dashboardId *int64) (result bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(deleteDashboardWidgets, *dashboardId, *userId)
	_, errW := pool.Exec(query)
	if errW != nil {
		log.Print(errW.Error(), errW)
		err = fmt.Errorf("unable to delete a dashboard widgets with specified dashboard id - %d", *dashboardId)
		return
	}
	result = true
	return
}

const updateDashboardWidget = "UPDATE amper.widget_sys SET label='%s', description='%s', configuration='%s' WHERE id=%d AND dashboard_id=%d AND created_by=%d"

func UpdateWidget(userId *int64, dashboardId *int64, widgetId *int64, label *string, description *string, configuration *string) (result bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(updateDashboardWidget, *label, *description, *configuration, *widgetId, *dashboardId, *userId)
	_, errW := pool.Exec(query)
	if errW != nil {
		log.Print(errW.Error(), errW)
		err = fmt.Errorf("unable to update a dashboard widget with specified parameters: label - %s, description - %s, configuration - %s", *label, *description, *configuration)
		return
	} /* else if count, _ := res.RowsAffected(); count < 1 {
		err = fmt.Errorf("unable to update a dashboard widget with specified parameters: label - %s, description - %s, configuration - %s", *label, *description, *configuration)
	}*/
	result = true
	return
}
