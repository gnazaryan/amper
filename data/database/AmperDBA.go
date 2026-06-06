package database

import (
	databasecache "amper/cache/database"
	"amper/common/structs"
	"amper/common/util"
	"amper/common/util/datetime"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

const getAmperInstanceQuery = "SELECT * FROM amper.amper_sys WHERE identifier=%d"

func GetInstance(identifier *int64, includeKey bool) (result *structs.Amper, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(getAmperInstanceQuery, *identifier)
	rows, errQ := pool.Query(query)
	if errQ == nil {
		for rows.Next() {
			var amper structs.Amper
			rows.Scan(&amper.Id, &amper.Identifier, &amper.Name, &amper.Type, &amper.Address, &amper.Port, &amper.State, &amper.StateUpdateDate, &amper.Usage, &amper.Limit, &amper.Directory, &amper.Key)
			if !includeKey {
				amper.Key = nil
			}
			result = &amper
			break
		}
	} else {
		err = errors.New("unable to run query against database to get amper instances")
		log.Print(errQ.Error(), errQ)
	}
	rows.Close()
	return result, err
}

const getAmperInstancesQuery = "SELECT * FROM amper.amper_sys"

func GetInstances(userID *int64, Type *string, includeKey bool) (instances []structs.Amper, err error) {
	var pool *sql.DB = databasecache.Pool()
	var query = getAmperInstancesQuery
	if Type != nil {
		query = query + " WHERE type='" + *Type + "'"
	}

	rows, errQ := pool.Query(query)
	if errQ == nil {
		for rows.Next() {
			var amper structs.Amper
			rows.Scan(&amper.Id, &amper.Identifier, &amper.Name, &amper.Type, &amper.Address, &amper.Port, &amper.State, &amper.StateUpdateDate, &amper.Usage, &amper.Limit, &amper.Directory, &amper.Key)
			if !includeKey {
				amper.Key = nil
			}
			instances = append(instances, amper)
		}
	} else {
		err = errors.New("unable to run query against database to get amper instances")
		log.Print(errQ.Error(), errQ)
	}
	rows.Close()
	return instances, err
}

const removeAmperInstance = "DELETE from amper.amper_sys where id=%d"

func RemoveInstance(userID *int64, amper structs.Amper) (bool, error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(removeAmperInstance, *amper.Id)
	_, errDB := pool.Exec(query)
	if errDB != nil {
		log.Println(errDB.Error(), errDB)
		return false, fmt.Errorf("unable to remove amper instance with id %d", amper.Id)
	}
	return true, nil
}

const insertAmperInstance = "INSERT INTO amper.amper_sys VALUES (null, %d, '%s', '%s', '%s', '%s', %d, '%s', %d, %d, '%s', '%s')"

func CreateInstance(userID *int64, amper structs.Amper) (bool, error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(insertAmperInstance, *amper.Identifier, *amper.Name, *amper.Type, *amper.Address, *amper.Port, *amper.State, datetime.GetDateTimeFormatted(), 0, *amper.Limit, *amper.Directory, *util.UUID())
	res, errW := pool.Exec(query)
	if errW != nil {
		log.Print(errW.Error(), errW)
		return false, fmt.Errorf("unable to insert an amper instance into the database with specified parameters: name - %s, address - %s, port - %s", *amper.Name, *amper.Address, *amper.Port)
	} else if count, _ := res.RowsAffected(); count < 1 {
		return false, fmt.Errorf("unable to insert an amper instance into the database with specified parameters: name - %s, address - %s, port - %s", *amper.Name, *amper.Address, *amper.Port)
	}
	return true, nil
}

const editAmperInstance = "UPDATE amper.amper_sys SET name='%s', address='%s', port='%s', limitation=%d, directory='%s' WHERE id=%d"

func EditInstance(userID *int64, amper structs.Amper) (success bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(editAmperInstance, *amper.Name, *amper.Address, *amper.Port, *amper.Limit, *amper.Directory, *amper.Id)
	_, errW := pool.Exec(query)
	if errW != nil {
		log.Print(errW.Error(), errW)
		return false, fmt.Errorf("unable to update amper instance with specified parameters: name - %s, address - %s, port - %s", *amper.Name, *amper.Address, *amper.Port)
	}
	return true, nil
}

const editLastUpdateDate = "UPDATE amper.amper_sys SET state_update_date='%s' WHERE identifier=%d"

func EditLastUpdateTime(userID *int64, identifier *int64, dateTime string) (success bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(editLastUpdateDate, dateTime, *identifier)
	_, errW := pool.Exec(query)
	if errW != nil {
		log.Print(errW.Error(), errW)
		return false, fmt.Errorf("unable to edit last update date of amper instance with specified parameters: identifier - %d", *identifier)
	}
	return true, nil
}
