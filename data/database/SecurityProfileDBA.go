package database

import (
	databasecache "amper/cache/database"
	"amper/common/structs"
	"amper/common/util"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

var profilesWithPagination = "SELECT P.id as profileId, P.name as profileName FROM amper.profile_sys as P %s %s LIMIT %d OFFSET %d"

// FetchProfiles query the database and retrieve the profiles for the provided start and limit parameters
func FetchProfiles(start *int, limit *int, search *[]string, sortField *string, sortDirection *string) (result []structs.Profile, err error) {
	var pool *sql.DB = databasecache.Pool()
	whereClous := ""
	if search != nil && len(*search) > 0 {
		whereClous = "WHERE (%s)"
		likes := make([]string, 0)
		searchColumns := [6]string{"P.name"}
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
	query := fmt.Sprintf(profilesWithPagination, whereClous, orderBy, *limit, *start)
	rows, errQ := pool.Query(query)
	if errQ == nil {
		for rows.Next() {
			var profile structs.Profile
			rows.Scan(&profile.ID, &profile.Name)
			result = append(result, profile)
		}
	} else {
		err = errors.New("unable to run query against database to get profiles")
		log.Print(errQ.Error(), errQ)
	}
	rows.Close()
	return
}
