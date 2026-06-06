package business

import (
	"amper/common/structs"
	"amper/data/database"
	"fmt"
	"log"
)

// FetchProfiles is responsible for retrieving users profiles with the provided start and limit parameters
func FetchProfiles(start *int, limit *int, search *[]string, sortField *string, sortDirection *string) (users []structs.Profile, err error) {
	users, errDb := database.FetchProfiles(start, limit, search, sortField, sortDirection)
	if errDb != nil {
		err = fmt.Errorf("unable to retrieve user profiles for the provided start: %d and limit: %d", start, limit)
		log.Print(errDb.Error(), errDb)
	}
	return
}
