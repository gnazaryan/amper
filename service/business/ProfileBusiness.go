package business

import (
	"amper/cache/business"
	"amper/common/crypto"
	"amper/common/structs"
	"amper/common/util"
	"amper/common/util/arrays"
	"amper/common/util/files"
	"amper/data/database"
	"amper/service/processor/rendition"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
)

const COVER_PHOTO_FILE = "Cover photo"
const PROFILE_PHOTO_FILE = "Profile photo"
const USER_DETAIL_FILE = "User detail"

var PROFILE_PATH = filepath.Join("__system__", "Profile")

// UpdateCover si designed to update the cover picture of the user profile
func UpdateCover(userId *int64, data []byte) (success bool, err error) {
	driveDir, errD := GetDriveDirectory(userId)
	if errD != nil {
		util.Loggify(errD)
		return false, fmt.Errorf("not able to locate the users active drive, contact the support")
	}
	directory := util.PointerString(filepath.Join(*driveDir, PROFILE_PATH))
	errA := files.CreateIfNotExists(*directory, 0755)
	if errA != nil {
		util.Loggify(errA)
		return false, fmt.Errorf("not able to initiate system directory in the users active drive, contact the support")
	}
	files, errorF := FetchFiles(userId, &PROFILE_PATH)
	if errorF != nil {
		util.Loggify(errorF)
		return false, fmt.Errorf("not able to look up system directory file in the users active drive, contact the support")
	}
	var existingFile *structs.FileMetadata
	for _, file := range *files {
		if !file.IsDir && *file.Name == COVER_PHOTO_FILE {
			existingFile = &file
			break
		}
	}
	size := int64(len(data))

	var errUP error
	if existingFile != nil {
		success, _, errUP = Upversion(userId, existingFile.Id, nil, &data, util.PointerString(COVER_PHOTO_FILE), util.PointerString("?"), &size, &PROFILE_PATH)
	} else {
		success, _, errUP = Upload(userId, nil, &data, util.PointerString(COVER_PHOTO_FILE), util.PointerString("?"), &size, &PROFILE_PATH, nil)
	}
	if errUP != nil {
		util.Loggify(errUP)
		err = fmt.Errorf("not able to update the cover photo of profile, contact the support")
	}
	return success, err
}

func ViewCover(userID *int64) (result *io.ReadCloser, metadata *structs.FileMetadata, err error) {
	files, errF := FetchFiles(userID, &PROFILE_PATH)
	var coverPhoto *structs.FileMetadata
	if errF == nil {
		for _, file := range *files {
			if *file.Name == COVER_PHOTO_FILE {
				coverPhoto = &file
				break
			}
		}
	}
	if coverPhoto != nil {
		return GetFile(userID, &PROFILE_PATH, coverPhoto.Id, coverPhoto.Version, nil)
	}
	return nil, nil, fmt.Errorf("we prefere you to have a cover photo uploaded first")
}

// UpdateCover si designed to update the cover picture of the user profile
func UpdatePhoto(userId *int64, data []byte) (success bool, err error) {
	driveDir, errD := GetDriveDirectory(userId)
	if errD != nil {
		util.Loggify(errD)
		return false, fmt.Errorf("not able to locate the users active drive, contact the support")
	}
	directory := util.PointerString(filepath.Join(*driveDir, PROFILE_PATH))
	errA := files.CreateIfNotExists(*directory, 0755)
	if errA != nil {
		util.Loggify(errA)
		return false, fmt.Errorf("not able to initiate system directory in the users active drive, contact the support")
	}
	files, errorF := FetchFiles(userId, &PROFILE_PATH)
	if errorF != nil {
		util.Loggify(errA)
		return false, fmt.Errorf("not able to look up system directory file in the users active drive, contact the support")
	}
	var existingFile *structs.FileMetadata
	for _, file := range *files {
		if !file.IsDir && *file.Name == PROFILE_PHOTO_FILE {
			existingFile = &file
		}
	}
	size := int64(len(data))

	var errUP error
	if existingFile != nil {
		success, _, errUP = Upversion(userId, existingFile.Id, nil, &data, util.PointerString(PROFILE_PHOTO_FILE), util.PointerString("?"), &size, &PROFILE_PATH)
	} else {
		success, _, errUP = Upload(userId, nil, &data, util.PointerString(PROFILE_PHOTO_FILE), util.PointerString("?"), &size, &PROFILE_PATH, nil)
	}
	if errUP != nil {
		util.Loggify(errUP)
		err = fmt.Errorf("not able to update the cover photo of profile, contact the support")
	}
	return success, err
}

func ViewPhoto(userID *int64) (result *io.ReadCloser, metadata *structs.FileMetadata, err error) {
	files, errF := FetchFiles(userID, &PROFILE_PATH)
	var profilePhoto *structs.FileMetadata
	if errF == nil {
		for _, file := range *files {
			if *file.Name == PROFILE_PHOTO_FILE {
				profilePhoto = &file
				break
			}
		}
	}
	if profilePhoto != nil {
		return GetFile(userID, &PROFILE_PATH, profilePhoto.Id, profilePhoto.Version, nil)
	}
	return nil, nil, fmt.Errorf("we prefere you to have a profile photo uploaded first")
}

func AdjustPhoto(userID *int64, sessionId *string, PositionX *int, PositionY *int, Width *int, Height *int) (success bool, result *string, err error) {
	files, errF := FetchFiles(userID, &PROFILE_PATH)
	var profilePhoto *structs.FileMetadata
	if errF == nil {
		for _, file := range *files {
			if *file.Name == PROFILE_PHOTO_FILE {
				profilePhoto = &file
				break
			}
		}
	} else {
		return false, nil, fmt.Errorf("we prefere you to have a profile photo uploaded first")
	}

	if profilePhoto != nil {
		reader, _, errR := GetFile(userID, util.PointerString(PROFILE_PATH), profilePhoto.Id, profilePhoto.Version, util.PointerBoolean(false))
		if errR != nil {
			util.Loggify(errR)
			return false, nil, fmt.Errorf("not able to retrieve the profile photo data, try again later")
		}
		photoData, errPD := io.ReadAll(*reader)
		if errPD != nil {
			util.Loggify(errPD)
			return false, nil, fmt.Errorf("not able to read the profile photo data, try again later")
		}
		crop, errC := rendition.Crop(&photoData, PositionX, PositionY, Width, Height)
		if errC == nil {
			cropBase64 := base64.StdEncoding.EncodeToString(crop)
			result = &cropBase64
			user, errU := database.GetUser(userID, nil, nil, true, nil, true, false)
			if errU == nil {
				errAP := user.AddConfig("profile", map[string]interface{}{
					"picture": map[string]interface{}{
						"PositionX": strconv.Itoa(*PositionX),
						"PositionY": strconv.Itoa(*PositionY),
						"Width":     strconv.Itoa(*Width),
						"Height":    strconv.Itoa(*Height),
					},
				})
				if errAP == nil {
					edited, errE := database.EditUserProperty(userID, map[string]string{
						"photo":  cropBase64,
						"config": *user.Config,
					})
					if errE != nil || !edited {
						util.Loggify(errE)
						return false, nil, fmt.Errorf("not able to update the profile photo")
					}
					instanceTypes := []string{"amperInstance"}
					successFC, errFC := FederatedCall(userID, sessionId, map[string]string{
						"amperDatastoreInstance": "amper-datastore/invalidateCache",
						"amperInstance":          "amper/invalidateCache",
					}, map[string]interface{}{
						"UserIdDelete": strconv.FormatInt(*user.ID, 10),
						"name":         "user",
					}, &instanceTypes)
					if !successFC || errFC != nil {
						util.Loggify(errFC)
						return false, nil, fmt.Errorf("AdjustPhoto - not able to reset the federated cache for user, contact the support")
					}
				}
			}
		}
	} else {
		return false, nil, fmt.Errorf("we prefere you to have a profile photo uploaded first")
	}
	return true, result, nil
}

func GetProfileState(userID *int64) (result map[string]interface{}, err error) {
	result = make(map[string]interface{})
	user := business.GetUser(userID, true)
	if user != nil {
		config, errC := user.GetConfig()
		if errC == nil {
			result["configuration"] = config
		} else {
			util.Loggify(errC)
			err = fmt.Errorf("not able retrieve the user configuration")
		}
	} else {
		err = fmt.Errorf("not able to locate the originating user")
	}
	files, errorF := FetchFiles(userID, &PROFILE_PATH)
	if errorF != nil {
		util.Loggify(errorF)
	}
	var existingFile *structs.FileMetadata
	if files != nil {
		for _, file := range *files {
			if !file.IsDir && *file.Name == USER_DETAIL_FILE {
				existingFile = &file
				break
			}
		}
	}
	if existingFile != nil {
		reader, errR := GetFileBody(userID, existingFile.Id)
		if errR != nil {
			util.Loggify(errR)
		} else {
			deteailBytes, errPD := io.ReadAll(*reader)
			if errPD != nil {
				util.Loggify(errPD)
			} else {
				var detail map[string]interface{}
				errUM := json.Unmarshal(deteailBytes, &detail)
				if errUM != nil {
					util.Loggify(errUM)
				} else {
					managers, okV := detail["managers"]
					if okV {
						var managerUsers []structs.User = make([]structs.User, 0)
						managersInterface := managers.([]interface{})
						managerIdsString, _ := arrays.InterfaceToString(&managersInterface)
						for _, managerIdString := range managerIdsString {
							managerId, errMI := strconv.ParseInt(managerIdString, 10, 64)
							if errMI == nil {
								managerUsers = append(managerUsers, *business.GetUser(&managerId, true))
							}
						}
						detail["managerUsers"] = managerUsers
					}
					reporters, okV := detail["reporters"]
					if okV {
						var reporterUsers []structs.User = make([]structs.User, 0)
						reportersInterface := reporters.([]interface{})
						reporterIdsString, _ := arrays.InterfaceToString(&reportersInterface)
						for _, reporterIdString := range reporterIdsString {
							reporterId, errRI := strconv.ParseInt(reporterIdString, 10, 64)
							if errRI == nil {
								reporterUsers = append(reporterUsers, *business.GetUser(&reporterId, true))
							}
						}
						detail["reporterUsers"] = reporterUsers
					}
					result["detail"] = detail
				}
			}
		}
	}
	return result, err
}

func SaveConfiguration(userId *int64, Name *string, Value *map[string]interface{}) (success bool, err error) {
	user, errU := database.GetUser(userId, nil, nil, true, nil, true, true)
	if errU == nil && user != nil {
		preprocessConfiguration(user, Name, Value)
		errAC := user.AddConfig(*Name, *Value)
		if errAC == nil {
			edited, errE := database.EditUserProperty(userId, map[string]string{
				"config": *user.Config,
			})
			if errE != nil || !edited {
				util.Loggify(errE)
				return false, fmt.Errorf("not able to update the profile configuration for %s", *Name)
			}
			return true, nil
		} else {
			util.Loggify(errAC)
			err = fmt.Errorf("not able to adjust the configuration to user's exisiting %s configuration ", *Name)
		}
	} else {
		util.Loggify(errU)
		err = fmt.Errorf("not able to retrieve a user to adjust the configuration for %s", *Name)
	}
	return false, err
}

func preprocessConfiguration(user *structs.User, Name *string, Value *map[string]interface{}) {
	switch *Name {
	case "settings":
		email, ok := (*Value)["email"].([]interface{})
		if ok {
			for _, item := range email {
				password, okP := item.(map[string]interface{})["password"].(string)
				if okP {
					item.(map[string]interface{})["password"] = hex.EncodeToString(crypto.Encrypt([]byte(password), *user.Password))
				}
				InitializeEmail(user.ID, item.(map[string]interface{}))
			}
		}
	}
}

func SaveDetail(userId *int64, Name *string, Value *string) (success bool, err error) {
	permittedFields := []string{"info", "about_me", "responsibilities", "skills"}
	if !arrays.ContainsS(permittedFields, *Name) {
		return false, fmt.Errorf("not able to update a user detail, since the field %s is not editable", *Name)
	}
	_, exists, errUD := database.GetUserDetail(userId)
	if !exists {
		_, errCUD := database.CreateUserDetail(userId)
		if errCUD != nil {
			util.Loggify(errCUD)
			return false, fmt.Errorf("not able to initialize a user detail for user %d", *userId)
		}
	}
	if errUD == nil {
		successU, errU := database.EditUserDetail(userId, Name, Value)
		if errU != nil || !successU {
			util.Loggify(errU)
			err = fmt.Errorf("not able to update the user detail for user id: %d, key: %s, value: %s", *userId, *Name, *Value)
		} else {
			success = true
			files, errorF := FetchFiles(userId, &PROFILE_PATH)
			if errorF != nil {
				util.Loggify(errorF)
			}
			var existingFile *structs.FileMetadata
			if files != nil {
				for _, file := range *files {
					if !file.IsDir && *file.Name == USER_DETAIL_FILE {
						existingFile = &file
						break
					}
				}
			}
			if existingFile == nil {
				detail := map[string]interface{}{
					*Name: *Value,
				}
				detailBytes, errUM := json.Marshal(detail)
				if errUM != nil {
					util.Loggify(errUM)
					return false, fmt.Errorf("SaveDetail - not able to convert user deatil to bytes, contact the support")
				}
				size := int64(len(detailBytes))
				success, _, errU := Upload(userId, nil, &detailBytes, util.PointerString(USER_DETAIL_FILE), util.PointerString("application/json"), &size, &PROFILE_PATH, nil)
				if !success || errU != nil {
					return false, fmt.Errorf("SaveDetail - not able to save the initial detail file, contact the support")
				}
			} else {
				reader, errR := GetFileBody(userId, existingFile.Id)
				if errR != nil {
					util.Loggify(errR)
					return false, fmt.Errorf("SaveDetail - not able to retrieve the detail file, contact the support")
				}
				deteailBytes, errPD := io.ReadAll(*reader)
				if errPD != nil {
					util.Loggify(errPD)
					return false, fmt.Errorf("SaveDetail - not able to read the detail file, contact the support")
				}
				var detail map[string]interface{}
				errUM := json.Unmarshal(deteailBytes, &detail)
				if errUM != nil {
					util.Loggify(errUM)
					return false, fmt.Errorf("SaveDetail - not able to unmarshal the detail file, contact the support")
				}
				detail[*Name] = *Value
				detailMarshaledBytes, errUM := json.Marshal(detail)
				if errUM != nil {
					util.Loggify(errUM)
					return false, fmt.Errorf("SaveDetail - not able to convert back user deatil to bytes, contact the support")
				}
				success, errU := Update(userId, existingFile.Id, &detailMarshaledBytes, util.PointerString(""))
				if errU != nil || !success {
					util.Loggify(errU)
					return false, fmt.Errorf("SaveDetail - not able to save the user detail file, contact the support")
				}
			}
		}
	}
	return success, err
}

func AddRelationship(userID *int64, employeeId int64, Type *string, Value int64) (success bool, err error) {
	permittedFields := []string{"manager", "reporter"}
	if !arrays.ContainsS(permittedFields, *Type) {
		return false, fmt.Errorf("not able to update the user detail file to add a %s", *Type)
	}
	files, errorF := FetchFiles(&employeeId, &PROFILE_PATH)
	if errorF != nil {
		util.Loggify(errorF)
	}
	var existingFile *structs.FileMetadata
	if files != nil {
		for _, file := range *files {
			if !file.IsDir && *file.Name == USER_DETAIL_FILE {
				existingFile = &file
				break
			}
		}
	}
	var key string
	if *Type == "reporter" {
		key = "reporters"
	} else if *Type == "manager" {
		key = "managers"
	}
	if existingFile == nil {
		detail := map[string]interface{}{
			key: []string{strconv.FormatInt(Value, 10)},
		}
		detailBytes, errUM := json.Marshal(detail)
		if errUM != nil {
			util.Loggify(errUM)
			return false, fmt.Errorf("AddRelationship - not able to convert user deatil to bytes, contact the support")
		}
		size := int64(len(detailBytes))
		success, _, errU := Upload(&employeeId, nil, &detailBytes, util.PointerString(USER_DETAIL_FILE), util.PointerString("application/json"), &size, &PROFILE_PATH, nil)
		if !success || errU != nil {
			return false, fmt.Errorf("AddRelationship - not able to save the initial detail file, contact the support")
		}
	} else {
		reader, errR := GetFileBody(&employeeId, existingFile.Id)
		if errR != nil {
			util.Loggify(errR)
			return false, fmt.Errorf("AddRelationship - not able to retrieve the detail file, contact the support")
		}
		deteailBytes, errPD := io.ReadAll(*reader)
		if errPD != nil {
			util.Loggify(errPD)
			return false, fmt.Errorf("AddRelationship - not able to read the detail file, contact the support")
		}
		var detail map[string]interface{}
		errUM := json.Unmarshal(deteailBytes, &detail)
		if errUM != nil {
			util.Loggify(errUM)
			return false, fmt.Errorf("AddRelationship - not able to unmarshal the detail file, contact the support")
		}
		values, okV := detail[key]
		if okV {
			valuesInterface := values.([]interface{})
			valuesString, _ := arrays.InterfaceToString(&valuesInterface)
			if !arrays.Contains(valuesString, strconv.FormatInt(Value, 10)) {
				valuesString = append(valuesString, strconv.FormatInt(Value, 10))
				detail[key] = valuesString
			}
		} else {
			detail[key] = []string{strconv.FormatInt(Value, 10)}
		}
		detailMarshaledBytes, errUM := json.Marshal(detail)
		if errUM != nil {
			util.Loggify(errUM)
			return false, fmt.Errorf("AddRelationship - not able to convert back user deatil to bytes, contact the support")
		}
		success, errU := Update(&employeeId, existingFile.Id, &detailMarshaledBytes, util.PointerString(""))
		if errU != nil || !success {
			util.Loggify(errU)
			return false, fmt.Errorf("AddRelationship - not able to save the user detail file, contact the support")
		}
	}
	return true, nil
}

func RemoveRelationship(userID *int64, employeeId int64, Type *string, Value int64) (success bool, err error) {
	permittedFields := []string{"manager", "reporter"}
	if !arrays.ContainsS(permittedFields, *Type) {
		return false, fmt.Errorf("not able to update the user detail file to add a %s", *Type)
	}
	files, errorF := FetchFiles(&employeeId, &PROFILE_PATH)
	if errorF != nil {
		util.Loggify(errorF)
	}
	var existingFile *structs.FileMetadata
	if files != nil {
		for _, file := range *files {
			if !file.IsDir && *file.Name == USER_DETAIL_FILE {
				existingFile = &file
				break
			}
		}
	}
	var key string
	if *Type == "reporter" {
		key = "reporters"
	} else if *Type == "manager" {
		key = "managers"
	}
	if existingFile != nil {
		reader, errR := GetFileBody(&employeeId, existingFile.Id)
		if errR != nil {
			util.Loggify(errR)
			return false, fmt.Errorf("RemoveRelationship - not able to retrieve the detail file, contact the support")
		}
		deteailBytes, errPD := io.ReadAll(*reader)
		if errPD != nil {
			util.Loggify(errPD)
			return false, fmt.Errorf("RemoveRelationship - not able to read the detail file, contact the support")
		}
		var detail map[string]interface{}
		errUM := json.Unmarshal(deteailBytes, &detail)
		if errUM != nil {
			util.Loggify(errUM)
			return false, fmt.Errorf("RemoveRelationship - not able to unmarshal the detail file, contact the support")
		}
		values, okV := detail[key]
		if okV {
			valuesInterface := values.([]interface{})
			valuesString, _ := arrays.InterfaceToString(&valuesInterface)
			if arrays.Contains(valuesString, strconv.FormatInt(Value, 10)) {
				valuesString = arrays.Remove(valuesString, strconv.FormatInt(Value, 10))
				detail[key] = valuesString
			}
		}
		detailMarshaledBytes, errUM := json.Marshal(detail)
		if errUM != nil {
			util.Loggify(errUM)
			return false, fmt.Errorf("RemoveRelationship - not able to convert back user deatil to bytes, contact the support")
		}
		success, errU := Update(&employeeId, existingFile.Id, &detailMarshaledBytes, util.PointerString(""))
		if errU != nil || !success {
			util.Loggify(errU)
			return false, fmt.Errorf("RemoveRelationship - not able to save the user detail file, contact the support")
		}
	}
	return true, nil
}
