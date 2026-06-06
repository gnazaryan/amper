package authorization

import (
	"amper/common/argument"
	"amper/common/structs"
	"amper/service/business"
)

func SaveSettings(userId *int64, settings *structs.Settings) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId, "settings": settings})
	if err != nil {
		return false, err
	}
	err = argument.Validate(map[string]interface{}{"settings.RootDirectory": settings.RootDirectory})
	if err != nil {
		return false, err
	}
	return business.SaveSettings(userId, settings)
}

func FetchSettings(userId *int64) (settings *structs.Settings, err error) {
	err = argument.Validate(map[string]interface{}{"userId": userId})
	if err != nil {
		return nil, err
	}
	return business.FetchSettings(userId)
}

func GetSetting(userId *int64, key *string, defaultValue *string) (setting string) {
	return business.GetSetting(userId, key, defaultValue)
}
