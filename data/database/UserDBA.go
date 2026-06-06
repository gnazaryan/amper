package database

import (
	databasecache "amper/cache/database"
	"amper/common/crypto"
	"amper/common/structs"
	"amper/common/util"
	"amper/common/util/ampstrings"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

const insertUserDetail string = "INSERT INTO amper.users_detail_sys VALUES (null, '%d', '%s', '%s', '%s', '%s')"

// CreateUser execute and create a user with the specified parameters
func CreateUserDetail(userId *int64) (result *int64, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(insertUserDetail, *userId, "{}", "", "", "{}")
	queryResult, errQ := pool.Exec(query)
	if errQ == nil {
		resultL, errL := queryResult.LastInsertId()
		if errL == nil {
			result = &resultL
		} else {
			err = fmt.Errorf("unable to create a user detail with the user id: %d", *userId)
			log.Print(errL.Error(), errL)
		}
	} else {
		err = fmt.Errorf("unable to execute create a user detail with the user id: %d", *userId)
		log.Print(errQ.Error(), errQ)
	}
	return result, err
}

const editUserDetail string = "UPDATE amper.users_detail_sys SET %s='%s' where user_id='%d'"

// EditUserDetail execute and update a user detail with the specified parameters
func EditUserDetail(userId *int64, Name *string, Value *string) (result bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(editUserDetail, *Name, *Value, *userId)
	queryResult, errQ := pool.Exec(query)
	if errQ == nil {
		_, errL := queryResult.RowsAffected()
		if errL == nil {
			result = true
		} else {
			err = fmt.Errorf("unable to update a user detail with the user id: %d", *userId)
			log.Print(errL.Error(), errL)
		}
	} else {
		err = fmt.Errorf("unable to execute update a user detail with the user id: %d", *userId)
		log.Print(errQ.Error(), errQ)
	}
	return result, err
}

const getUserDetail string = "SELECT * FROM amper.users_detail_sys where user_id='%d'"

func GetUserDetail(id *int64) (result *structs.UserDetail, exists bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(getUserDetail, *id)
	var row *sql.Row = pool.QueryRow(query)
	if row.Err() == nil {
		result = new(structs.UserDetail)
		errS := row.Scan(&result.ID, &result.UserId, &result.Info, &result.AboutMe, &result.Responsibilities, &result.Skills)
		if errS == sql.ErrNoRows {
			return nil, false, nil
		} else if errS != nil {
			util.Loggify(errS)
			err = fmt.Errorf("not able to retrieve the user detail for user: %d", *id)
		}
	} else {
		util.Loggify(row.Err())
		err = fmt.Errorf("not able to retrieve the user detail for user: %d", *id)
	}
	return result, true, err
}

const getUserWithActivationID string = "SELECT  U.id as id, U.firstName as firstName, U.lastName as lastName, U.middleName as middleName, U.username as username, U.password as password, U.photo as photo, U.profile as profile, U.email as email, U.active as active, U.amperId as amperId, U.state as state, U.activationCode as activationCode, U.config as config FROM amper.users_sys as U where %s"

// GetUser returns database row for the provided username
func GetUser(id *int64, username *string, activationCode *string, active bool, deleted *bool, config bool, includePassword bool) (result *structs.User, err error) {
	var pool *sql.DB = databasecache.Pool()
	whereClouses := make([]string, 0)
	if active {
		whereClouses = append(whereClouses, fmt.Sprintf("active=%d", 1))
	}
	if deleted != nil && *deleted {
		whereClouses = append(whereClouses, fmt.Sprintf("U.state in (%d) ", 0))
	} else if deleted != nil && !*deleted {
		whereClouses = append(whereClouses, fmt.Sprintf("U.state in (%d) ", 1))
	}
	if id != nil {
		whereClouses = append(whereClouses, fmt.Sprintf("id=%d", *id))
	} else if username != nil {
		whereClouses = append(whereClouses, fmt.Sprintf("username='%s'", *username))
	} else if activationCode != nil {
		whereClouses = append(whereClouses, fmt.Sprintf("activationCode='%s'", *activationCode))
	}
	query := fmt.Sprintf(getUserWithActivationID, strings.Join(whereClouses, " AND "))
	var row *sql.Row = pool.QueryRow(query)
	var errQ error
	if row.Err() == nil {
		result = new(structs.User)
		errQ = row.Scan(&result.ID, &result.FirstName, &result.LastName, &result.MiddleName, &result.Username, &result.Password,
			&result.Photo, &result.Profile, &result.Email, &result.Active, &result.AmperId, &result.State, &result.ActivationCode, &result.Config)
		if !config {
			result.Config = nil
		}
		if !includePassword {
			result.Password = nil
		}
	}
	if errQ != nil {
		result = nil
		err = errors.New("unable to execute query to get user")
		log.Print(errQ.Error(), errQ)
	}
	return
}

const getUsersWithPagination = "SELECT U.id as id, U.firstName as firstName, U.middleName as middleName, U.lastName as lastName, U.username as username, U.photo as photo, U.email as email, U.active as active, U.amperId as amperId, U.state as state, P.id as profileId, P.name as profileName FROM amper.users_sys as U INNER JOIN amper.profile_sys as P on U.profile = P.id %s %s LIMIT %d OFFSET %d"
const getUsersWithPaginationCount = "SELECT Count(U.id) as totalCount FROM amper.users_sys as U %s"

// GetUsers query the database and retrieve the users for the provided start and limit parameters
func GetUsers(start *int, limit *int, search *[]string, sortField *string, sortDirection *string) (result []structs.UserAndProfile, totalCount int, err error) {
	var pool *sql.DB = databasecache.Pool()
	whereClous := "WHERE U.state in (1)"
	if search != nil && len(*search) > 0 {
		whereClous = whereClous + " AND (%s)"
		likes := make([]string, 0)
		searchColumns := [5]string{"U.firstName", "U.middleName", "U.lastName", "U.username", "U.email"}
		for _, value := range searchColumns {
			likes = append(likes, fmt.Sprintf("%s LIKE '%%%s%%'", value, strings.Join(*search, " ")))
			for _, searchValue := range *search {
				like := fmt.Sprintf("%s LIKE '%%%s%%'", value, searchValue)
				likes = append(likes, like)
			}
		}
		whereClous = fmt.Sprintf(whereClous, strings.Join(likes, " OR "))
	}
	orderBy := ""
	if sortField != nil && sortDirection != nil {
		direction := util.IfElse(strings.ToLower(*sortDirection) == "asc", "ASC", "DESC")
		orderBy = fmt.Sprintf("ORDER BY %s %s", *sortField, direction)
	}
	query := fmt.Sprintf(getUsersWithPagination, whereClous, orderBy, *limit, *start)
	rows, errQ := pool.Query(query)
	if errQ == nil {
		for rows.Next() {
			var userAndProfile structs.UserAndProfile
			rows.Scan(&userAndProfile.ID, &userAndProfile.FirstName, &userAndProfile.MiddleName, &userAndProfile.LastName, &userAndProfile.Username,
				&userAndProfile.Photo, &userAndProfile.Email, &userAndProfile.Active, &userAndProfile.AmperId, &userAndProfile.State, &userAndProfile.ProfileID, &userAndProfile.ProfileName)
			result = append(result, userAndProfile)
		}
	} else {
		err = errors.New("unable to run query against database to get users")
		log.Print(errQ.Error(), errQ)
	}
	rows.Close()

	totalCountQuery := fmt.Sprintf(getUsersWithPaginationCount, whereClous)
	row := pool.QueryRow(totalCountQuery)
	if row.Err() == nil {
		row.Scan(&totalCount)
	} else {
		log.Println(row.Err(), row.Err())
		//skip reporting the totalc count error
	}
	return
}

const getUsersInWithPagination = "SELECT U.id as id, U.firstName as firstName, U.middleName as middleName, U.lastName as lastName, U.username as username, U.photo as photo, U.email as email, U.active as active, U.amperId as amperId, U.state as state FROM amper.users_sys as U %s LIMIT %d OFFSET %d"
const getUsersInWithPaginationCount = "SELECT Count(U.id) as totalCount FROM amper.users_sys as U %s"

// GetUsers query the database and retrieve the users for the provided start and limit parameters
func GetUsersIn(start *int, limit *int, search *[]string, userIds []string) (result []structs.User, resultTotalCount int, err error) {
	var pool *sql.DB = databasecache.Pool()
	whereClous := "WHERE U.id in (" + strings.Join(userIds[:], ",") + ")"
	if search != nil && len(*search) > 0 {
		whereClous = whereClous + " AND (%s)"
		likes := make([]string, 0)
		searchColumns := [5]string{"U.firstName", "U.middleName", "U.lastName", "U.username", "U.email"}
		for _, value := range searchColumns {
			likes = append(likes, fmt.Sprintf("%s LIKE '%%%s%%'", value, strings.Join(*search, " ")))
			for _, searchValue := range *search {
				like := fmt.Sprintf("%s LIKE '%%%s%%'", value, searchValue)
				likes = append(likes, like)
			}
		}
		whereClous = fmt.Sprintf(whereClous, strings.Join(likes, " OR "))
	}

	query := fmt.Sprintf(getUsersInWithPagination, whereClous, *limit, *start)
	rows, errQ := pool.Query(query)
	if errQ == nil {
		for rows.Next() {
			var user structs.User
			rows.Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.Username,
				&user.Photo, &user.Email, &user.Active, &user.AmperId, &user.State)
			result = append(result, user)
		}
	} else {
		err = errors.New("unable to run query against database to get users")
		log.Print(errQ.Error(), errQ)
	}
	rows.Close()

	totalCountQuery := fmt.Sprintf(getUsersInWithPaginationCount, whereClous)
	row := pool.QueryRow(totalCountQuery)
	if row.Err() == nil {
		row.Scan(&resultTotalCount)
	} else {
		log.Println(row.Err(), row.Err())
		//skip reporting the totalc count error
	}
	return
}

const insertUser string = "INSERT INTO amper.users_sys VALUES (null, '%s', '%s', '%s', '%s', '', '%s', '%d', '%s', '%d', '%s', '%d', '1', null)"

var passive int = 0

// CreateUser execute and create a user with the specified parameters
func CreateUser(user structs.User) (result *int64, err error) {
	var pool *sql.DB = databasecache.Pool()
	user.Active = &passive
	user.ActivationCode = util.PointerString(crypto.UUID())
	query := fmt.Sprintf(insertUser, *user.FirstName, *user.LastName, ampstrings.EmptyIfNil(user.MiddleName), *user.Username, ampstrings.EmptyIfNil(user.MiddleName), *user.Profile, *user.Email, *user.Active, *user.ActivationCode, *user.AmperId)
	queryResult, errQ := pool.Exec(query)
	if errQ == nil {
		resultL, errL := queryResult.LastInsertId()
		if errL == nil {
			result = &resultL
		} else {
			err = fmt.Errorf("unable to create a user with the username: %s", *user.Username)
			log.Print(errL.Error(), errL)
		}
	} else {
		err = fmt.Errorf("unable to create a user with the username: %s", *user.Username)
		log.Print(errQ.Error(), errQ)
	}
	return
}

const editUser string = "UPDATE amper.users_sys SET firstName='%s', lastName='%s', middleName='%s', photo='%s', profile='%d', email='%s' where id='%d'"

// EditUser execute and update a user with the specified parameters
func EditUser(user structs.User) (result bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(editUser, *user.FirstName, *user.LastName, *user.MiddleName, *user.Photo, *user.Profile, *user.Email, *user.ID)
	queryResult, errQ := pool.Exec(query)
	if errQ == nil {
		_, errL := queryResult.RowsAffected()
		if errL == nil {
			result = true
		} else {
			err = fmt.Errorf("unable to update a user with the username: %s", *user.Username)
			if errL != nil {
				log.Print(errL.Error(), errL)
			}
		}
	} else {
		err = fmt.Errorf("unable to update a user with the username: %s", *user.Username)
		if errQ != nil {
			log.Print(errQ.Error(), errQ)
		}
	}
	return result, err
}

const editUserByProperty string = "UPDATE amper.users_sys SET %s where id='%d'"

// EditUser execute and update a user with the specified parameter propert values
func EditUserProperty(userID *int64, values map[string]string) (result bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	var propValues []string
	for key, value := range values {
		propValues = append(propValues, fmt.Sprintf("%s='%s'", key, value))
	}

	query := fmt.Sprintf(editUserByProperty, strings.Join(propValues, ","), *userID)
	queryResult, errQ := pool.Exec(query)
	if errQ == nil {
		rowsAffected, errL := queryResult.RowsAffected()
		if errL == nil && rowsAffected > 0 {
			result = true
		} else {
			err = fmt.Errorf("unable to update a user with the user id: %d", *userID)
			if errL != nil {
				log.Print(errL.Error(), errL)
			}
		}
	} else {
		err = fmt.Errorf("unable to update a user with the user id: %d", *userID)
		if errQ != nil {
			log.Print(errQ.Error(), errQ)
		}
	}
	return
}

const activateUser string = "UPDATE amper.users_sys SET active=1, activationCode='%s', password='%s' where id='%d'"

// Activate execute and activate a user with the give activation code
func Activate(userID *int64, activationCode *string, password *string) (result bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(activateUser, *activationCode, *password, *userID)
	queryResult, errQ := pool.Exec(query)
	if errQ == nil {
		resultL, errL := queryResult.RowsAffected()
		if errL == nil && resultL == 1 {
			result = true
		} else {
			err = fmt.Errorf("unable to activate a user with id: %d", *userID)
			log.Print(errL.Error(), errL)
		}
	} else {
		err = fmt.Errorf("unable to update a user with the username: %d", *userID)
		log.Print(errQ.Error(), errQ)
	}
	return
}

var deleteUser string = "DELETE FROM amper.users_sys where id=%d"

// Remove executes and deletes a user ith the given userId
func Remove(userIDToRemove *int64) (result bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(deleteUser, *userIDToRemove)
	queryResult, errQ := pool.Exec(query)
	if errQ == nil {
		resultL, errL := queryResult.RowsAffected()
		if errL == nil && resultL == 1 {
			result = true
		} else {
			err = fmt.Errorf("unable to delete a user with id: %d", *userIDToRemove)
			log.Print(errL.Error(), errL)
		}
	} else {
		err = fmt.Errorf("unable to delete a user with id: %d", *userIDToRemove)
		log.Print(errQ.Error(), errQ)
	}
	return
}

var deleteUserSoft string = "UPDATE amper.users_sys SET state='0', active='0' where id='%d'"

// Remove executes and deletes a user ith the given userId
func RemoveSoft(userIDToRemove *int64) (result bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(deleteUserSoft, *userIDToRemove)
	queryResult, errQ := pool.Exec(query)
	if errQ == nil {
		resultL, errL := queryResult.RowsAffected()
		if errL == nil && resultL == 1 {
			result = true
		} else {
			err = fmt.Errorf("unable to delete a user with id: %d", *userIDToRemove)
			log.Print(errL.Error(), errL)
		}
	} else {
		err = fmt.Errorf("unable to delete a user with id: %d", *userIDToRemove)
		log.Print(errQ.Error(), errQ)
	}
	return
}
