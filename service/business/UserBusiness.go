package business

import (
	"amper/cache/business"
	"amper/common/crypto"
	"amper/common/notification"
	"amper/common/structs"
	"amper/common/util"
	"amper/common/util/ampstrings"
	"amper/common/util/datetime"
	"amper/data/database"
	"amper/properties/application"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var SYSTEM_USER_ID int64 = 1

// UserLogin is responsibel for applying a business logic to the user authentication process
func UserLogin(username *string, password *string) (result *structs.User, err error) {
	user, errDb := database.GetUser(nil, username, nil, true, util.PointerBoolean(false), true, true)
	if user != nil && *user.Active == 1 {
		user.Initialize(true)
		passwordEncrypted, _ := hex.DecodeString(*user.Password)
		passwordDecrypted := string(crypto.Decrypt(passwordEncrypted, crypto.Passphrase))
		if passwordDecrypted == *password {
			result = user
			result.SessionID = GenerateSessionId(user.ID, user.Password, util.PointerString("user"))
			user.Password = nil
			user.Config = nil
		} else {
			err = fmt.Errorf("unable to authenticate the user with the provided username: %s and password: ?", *username)
			if errDb != nil {
				log.Print(errDb.Error(), errDb)
			}
		}
	} else {
		if user == nil {
			err = fmt.Errorf("no user was found with the username: %s provided", *username)
		} else if *user.Active != 1 {
			err = fmt.Errorf("the user with username: %s is not active, please follow the email directions and activate", *username)
		}

		if errDb != nil {
			log.Print(errDb.Error(), errDb)
		}
	}
	return
}

// The function is used to generate session id
// !!!ATTENTION!!! The same copy of code is used inside the datastore repository, if this method is changed or fixed
// Consider making the changes or fixing the issues in the datastore repo too
func GenerateSessionId(id *int64, passphrase *string, Type *string) string {
	result := fmt.Sprintf("%s%s%d", datetime.GetDateTimeFormatted(), ampstrings.SEPERATOR, *id)
	result = string(crypto.Encrypt([]byte(result), *passphrase))
	result = fmt.Sprintf("%s%s%d%s%s", result, ampstrings.SEPERATOR, *id, ampstrings.SEPERATOR, *Type)
	result = string(crypto.Encrypt([]byte(result), crypto.Passphrase))
	result = base64.StdEncoding.EncodeToString([]byte(result))
	return result
}

var sessionDuration time.Duration = 3 * time.Hour

func ValidateSession(sessionId *string) (bool, int64, *string) {
	if sessionId != nil {
		sessionIdDecoded, _ := base64.StdEncoding.DecodeString(*sessionId)
		sessionIdDecodedDecrypted := crypto.Decrypt([]byte(sessionIdDecoded), crypto.Passphrase)
		sessionSplit := strings.Split(string(sessionIdDecodedDecrypted), ampstrings.SEPERATOR)
		if len(sessionSplit) == 3 {
			id, errU := strconv.ParseInt(sessionSplit[1], 10, 64)
			if errU == nil {
				var passphrase *string
				if sessionSplit[2] == "user" {
					user := business.GetUser(&id, false)
					passphrase = user.Password
				} else {
					instance := business.GetAmperInstance(id)
					passphrase = instance.Key
				}
				if passphrase != nil {
					sessionIdDecrypted := string(crypto.Decrypt([]byte(sessionSplit[0]), *passphrase))
					sessionIdDecryptedSplit := strings.Split(sessionIdDecrypted, ampstrings.SEPERATOR)
					if len(sessionIdDecryptedSplit) == 2 {
						idInternal, errU := strconv.ParseInt(sessionIdDecryptedSplit[1], 10, 64)
						if errU == nil && idInternal == id {
							sessionTime, errST := datetime.ParseDateTime(util.PointerString(sessionIdDecryptedSplit[0]))
							if errST == nil {
								now, _ := datetime.ParseDateTime(util.PointerString(datetime.GetDateTimeFormatted()))
								if now.Before(sessionTime.Add(sessionDuration)) {
									return true, id, sessionId
								} else {
									util.Loggify(fmt.Errorf("session has already been expired"))
								}
							} else {
								util.Loggify(fmt.Errorf("session id doesn't contain a valid encrypted date time"))
							}
						} else {
							util.Loggify(fmt.Errorf("session id doesn't contain a valid user id"))
						}
					} else {
						util.Loggify(fmt.Errorf("session id is not in a valid form"))
					}
				} else {
					util.Loggify(fmt.Errorf("not able to locate the user to authenticate"))
				}
			} else {
				util.Loggify(fmt.Errorf("not able to locate the user id inside the session id"))
			}
		}
	}
	return false, -1, nil
}

// GetUsers is responsible for retrieving users with the provided start and limit parameters
func GetUsers(start *int, limit *int, search *[]string, sortField *string, sortDirection *string) (users []structs.UserAndProfile, totalCount int, err error) {
	users, totalCount, errDb := database.GetUsers(start, limit, search, sortField, sortDirection)
	if errDb != nil {
		err = fmt.Errorf("unable to retrieve users for the provided start: %d and limit: %d", start, limit)
		log.Print(errDb.Error(), errDb)
	}
	return
}

// CreateUser is responsable for creating a user with the provided parameters and
// checking if there is already a user with same username
func CreateUser(userID *int64, user structs.User) (result bool, err error) {
	userDb, errDb := database.GetUser(nil, user.Username, nil, false, nil, false, false)
	if userDb == nil {
		resultDb, errDb1 := database.CreateUser(user)
		if resultDb == nil || errDb1 != nil {
			err = fmt.Errorf("unable to create a user with the provided username: %s", *user.Username)
		} else {
			userDb, errDb = database.GetUser(resultDb, nil, nil, false, nil, false, false)
			result = true
			if userDb == nil || errDb != nil {
				err = fmt.Errorf("unable to create a user with the provided username: %s", *user.Username)
				result = false
			}
		}
	} else {
		err = fmt.Errorf("another user already exists with provided username: %s, we prefere non repeting usernames", *user.Username)
		if errDb != nil {
			log.Print(errDb.Error(), errDb)
		}
	}
	if result && err == nil {
		config, errAP := application.Get()
		if errAP != nil {
			log.Printf("unable to get application properties: %v", errAP)
		}
		uinodes := config.GetString("uinodes", "http://dev.amper.cloud:3000/")
		uinode := strings.Split(uinodes, ",")[0]

		resultN, errN := notification.Send(userDb, util.PointerString("userRegistration"), structs.UserNotification{
			UserFirstName: userDb.FirstName,
			ButtonLabel:   util.PointerString("Verify and set password"),
			ButtonHref:    util.PointerString(uinode + "?activationCode=" + *userDb.ActivationCode),
		})

		if errN != nil || resultN {
			if errN != nil {
				log.Printf("unable to send notification for user registration with username: %s and error : %v", *user.Username, errN)
			}
			err = fmt.Errorf("not able to send registration notification for user %s", *user.Username)
		}
	}
	return
}

// EditUser is responsable for updating/modifying an existing user with the provided
// new parameters
func EditUser(userID *int64, sessionId *string, user structs.User) (result bool, err error) {
	userDb, errDb := database.GetUser(nil, user.Username, nil, false, nil, false, false)
	if userDb != nil {
		success, errDb1 := database.EditUser(user)
		if !success || errDb1 != nil {
			err = fmt.Errorf("unable to update a user with the provided username: %s", *user.Username)
		} else {
			result = true
			instanceTypes := []string{"amperInstance"}
			successFC, errFC := FederatedCall(userID, sessionId, map[string]string{
				"amperDatastoreInstance": "amper-datastore/invalidateCache",
				"amperInstance":          "amper/invalidateCache",
			}, map[string]interface{}{
				"UserIdDelete": strconv.FormatInt(*user.ID, 10),
				"name":         "user",
			}, &instanceTypes)
			util.Loggify(errFC)
			if !successFC || errFC != nil {
				err = fmt.Errorf("not able to reset the federated cache for user, contact the support")
				result = false
			}
		}
	} else {
		err = fmt.Errorf("a user with the provided id: %d does not exist", *user.ID)
		if errDb != nil {
			log.Print(errDb.Error(), errDb)
		}
	}
	return
}

// Activate is responsible for activating a user with
func Activate(activationCode *string, password *string) (result bool, err error) {
	userDb, errDb := database.GetUser(nil, nil, activationCode, false, util.PointerBoolean(false), false, false)
	if userDb != nil && *userDb.Active == 0 {
		passwordEncypted := crypto.Encrypt([]byte(*password), crypto.Passphrase)
		passwordEncryptedEncoded := hex.EncodeToString(passwordEncypted)
		success, errDb1 := database.Activate(userDb.ID, activationCode, &passwordEncryptedEncoded)
		if !success || errDb1 != nil {
			err = fmt.Errorf("unable to activate a user with the provided username: %s", *userDb.Username)
		} else {
			result = true
		}
	} else {
		err = fmt.Errorf("a user with the provided activation code: %s does not exist", *activationCode)
		if errDb != nil {
			log.Print(errDb.Error(), errDb)
		}
	}
	return
}

// IsValidUserName is responsible for checking if the given username is available
func IsValidUserName(userID *int64, username *string) (result bool, err error) {
	userDb, errDb := database.GetUser(nil, username, nil, false, nil, false, false)
	if userDb != nil {
		result = true
	} else if errDb != nil {
		err = fmt.Errorf("unable to check if a user with the provided username: %s exist", *username)
		log.Print(errDb.Error(), errDb)
	}
	return
}

// Remove is responsible for removing a user with the provided id if exists
func Remove(userID *int64, userIDToRemove *int64) (result bool, err error) {
	userDb, errDb := database.GetUser(userIDToRemove, nil, nil, false, nil, false, false)
	if userDb != nil {
		success, errDb := database.Remove(userIDToRemove)
		if !success || errDb != nil {
			err = fmt.Errorf("unable to delete a user with the provided id: %d", *userIDToRemove)
		} else {
			result = true
		}
	} else if errDb != nil {
		err = fmt.Errorf("a user with the provided user id: %d does not exist", *userID)
		log.Print(errDb.Error(), errDb)
	}
	return
}

// Remove is responsible for removing a user with the provided id if exists
func RemoveSoft(userID *int64, userIDToRemove *int64) (result bool, err error) {
	userDb, errDb := database.GetUser(userIDToRemove, nil, nil, false, util.PointerBoolean(false), false, false)
	if userDb != nil {
		success, errDb := database.RemoveSoft(userIDToRemove)
		if !success || errDb != nil {
			err = fmt.Errorf("unable to delete a user with the provided id: %d", *userIDToRemove)
		} else {
			result = true
		}
	} else if errDb != nil {
		err = fmt.Errorf("a user with the provided user id: %d does not exist", *userID)
		log.Print(errDb.Error(), errDb)
	}
	return
}
