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

const getDashboard = "SELECT * FROM amper.dashboard_sys WHERE created_by=%d AND id=%d"

func GetDashboard(userId *int64, dashboardId *int64) (result structs.Dashboard, err error) {
	var pool *sql.DB = databasecache.Pool()
	rows, errQ := pool.Query(fmt.Sprintf(getDashboard, *userId, *dashboardId))
	if errQ == nil {
		for rows.Next() {
			rows.Scan(&result.ID, &result.Label, &result.Description, &result.Configuration, &result.CreatedDate, &result.CreatedBy)
			break
		}
	} else {
		err = errors.New("unable to run query against database to get dashboard")
		log.Print(errQ.Error(), errQ)
	}
	rows.Close()
	return
}

const getDashboards = "SELECT * FROM amper.dashboard_sys WHERE created_by=%d"

func GetDashboards(userId *int64) (result *[]structs.Dashboard, err error) {
	var pool *sql.DB = databasecache.Pool()
	rows, errQ := pool.Query(fmt.Sprintf(getDashboards, *userId))
	dashboards := []structs.Dashboard{}
	if errQ == nil {
		for rows.Next() {
			var dashboard structs.Dashboard
			rows.Scan(&dashboard.ID, &dashboard.Label, &dashboard.Description, &dashboard.Configuration, &dashboard.CreatedDate, &dashboard.CreatedBy)
			dashboards = append(dashboards, dashboard)
		}
	} else {
		err = errors.New("unable to run query against database to get dashboards")
		log.Print(errQ.Error(), errQ)
	}
	rows.Close()
	result = &dashboards
	return
}

const insertDashboard = "INSERT INTO amper.dashboard_sys VALUES (null, '%s', '%s', '%s', '%s', %d)"

func AddDashboard(userId *int64, label *string, description *string, configuration *string) (result bool, errRes error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(insertDashboard, *label, *description, *configuration, datetime.GetDateTimeFormatted(), *userId)
	res, err := pool.Exec(query)
	if err != nil {
		log.Print(err.Error(), err)
		errRes = fmt.Errorf("unable to inser a dashboard with specified parameters: label - %s, description - %s, configuration - %s", *label, *description, *configuration)
		return
	} else if count, _ := res.RowsAffected(); count < 1 {
		errRes = fmt.Errorf("unable to inser a dashboard with specified parameters: label - %s, description - %s, configuration - %s", *label, *description, *configuration)
	}
	result = true
	return
}

const updateDashboard = "UPDATE amper.dashboard_sys SET label='%s', description='%s', configuration='%s' WHERE id=%d AND created_by=%d";
func UpdateDashboard(userId *int64, id *int64, label *string, description *string, configuration *string) (result bool, errRes error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(updateDashboard, *label, *description, *configuration, *id, *userId)
	res, err := pool.Exec(query)
	if err != nil {
		log.Print(err.Error(), err)
		errRes = fmt.Errorf("unable to update a dashboard with specified parameters: label - %s, description - %s, configuration - %s", *label, *description, *configuration)
		return
	} else if count, _ := res.RowsAffected(); count < 1 {
		errRes = fmt.Errorf("unable to update a dashboard with specified parameters: label - %s, description - %s, configuration - %s", *label, *description, *configuration)
	}
	result = true
	return
}

const deleteDashboard = "DELETE from amper.dashboard_sys where id=%d AND created_by=%d"

func RemoveDashboard(userId *int64, dashboardId *int64) (result bool, errRes error) {
	success, errW := RemoveWidgets(userId, dashboardId)
	if (errW == nil && success) {
		var pool *sql.DB = databasecache.Pool()
		query := fmt.Sprintf(deleteDashboard, *dashboardId, *userId)
		res, err := pool.Exec(query)
		if err != nil {
			log.Print(err.Error(), err)
			errRes = fmt.Errorf("unable to delete a dashboard with specified id - %d", *dashboardId)
			result = false
			return
		} else if count, _ := res.RowsAffected(); count < 1 {
			errRes = fmt.Errorf("unable to delete a dashboard with specified id - %d", *dashboardId)
			result = false
			return
		}
		result = true
	} else {
		result = false
		if (errW != nil) {
			log.Print(errW.Error(), errW)
		}
		errRes = fmt.Errorf("unable to delete a dashboard with specified id - %d", *dashboardId)
	}
	return
}
