package database

import (
	databasecache "amper/cache/database"
	"amper/common/structs"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

const getUserRelationships string = "SELECT * FROM amper.users_relationship_sys WHERE employee_id=%d"

func GetUserRelationships(userId int64, emloyeeId int64) (result []structs.UserRelationship, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(getUserRelationships, emloyeeId)
	rows, errQ := pool.Query(query)
	if errQ == nil {
		for rows.Next() {
			var userRelationship structs.UserRelationship
			rows.Scan(&userRelationship.ID, &userRelationship.EmployeeId, &userRelationship.ManagerId)
			result = append(result, userRelationship)
		}
	} else {
		err = errors.New("unable to run query against database to get user relationships")
		log.Print(errQ.Error(), errQ)
	}
	rows.Close()
	return result, err
}

const deleteUserRelationship string = "DELETE FROM amper.users_relationship_sys WHERE manager_id=%d AND employee_id=%d"

func DeleteUserRelationship(userId int64, managerId int64, employeeId int64) (result bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(deleteUserRelationship, managerId, employeeId)
	queryResult, errQ := pool.Exec(query)
	if errQ == nil {
		_, errL := queryResult.RowsAffected()
		if errL == nil {
			result = true
		} else {
			err = fmt.Errorf("unable to delete a user relationship with the manager id: %d", managerId)
			log.Print(errL.Error(), errL)
		}
	} else {
		err = fmt.Errorf("unable to execute delete a user relationship with the manager id: %d", managerId)
		log.Print(errQ.Error(), errQ)
	}
	return result, err
}

const insertUserRelationship string = "INSERT INTO amper.users_relationship_sys VALUES (null, '%d', '%d')"

// CreateUser execute and create a user with the specified parameters
func CreateUserRelationship(userId int64, emloyeeId int64, managerId int64) (result *int64, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(insertUserRelationship, emloyeeId, managerId)
	queryResult, errQ := pool.Exec(query)
	if errQ == nil {
		resultL, errL := queryResult.LastInsertId()
		if errL == nil {
			result = &resultL
		} else {
			err = fmt.Errorf("unable to create a user relationship with the user id: %d", emloyeeId)
			log.Print(errL.Error(), errL)
		}
	} else {
		err = fmt.Errorf("unable to execute create a user relationship with the user id: %d", emloyeeId)
		log.Print(errQ.Error(), errQ)
	}
	return result, err
}
