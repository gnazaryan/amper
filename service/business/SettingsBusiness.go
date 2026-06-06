package business

import (
	"amper/common/structs"
	"amper/common/util"
	"amper/common/util/ampstrings"
	"amper/properties/application"
	"encoding/json"
	"fmt"
	"log"
)

func SaveSettings(userId *int64, settings *structs.Settings) (success bool, err error) {
	imapB, err := json.Marshal(settings.Imap)
	smtpB, err := json.Marshal(settings.Smtp)
	return application.SetValue(
		map[string]interface{}{
			"amper.rootDirectory":   *settings.RootDirectory,
			"amper.adobeLicenseKey": ampstrings.EmptyIfNil(settings.AdobeLicenseKey),
			"amper.imap":            string(imapB),
			"amper.smtp":            string(smtpB),
		})
}

func FetchSettings(userId *int64) (settings *structs.Settings, err error) {
	appConfiguration, errAC := application.Get()
	if errAC != nil {
		log.Println(errAC.Error(), errAC)
		return nil, fmt.Errorf("not able to load th eapplication configuration, try again or contact the support")
	}
	imapS := appConfiguration.GetString("amper.imap", "")
	imap := structs.Imap{
		Domains: make([]structs.Domain, 0),
	}
	json.Unmarshal([]byte(imapS), &imap)

	smtpS := appConfiguration.GetString("amper.smtp", "")
	smtp := structs.Smtp{
		Domains: make([]structs.Domain, 0),
	}
	json.Unmarshal([]byte(smtpS), &smtp)

	result := structs.Settings{
		RootDirectory:   util.PointerString(appConfiguration.GetString("amper.rootDirectory", "")),
		AdobeLicenseKey: util.PointerString(appConfiguration.GetString("amper.adobeLicenseKey", "")),
		Imap:            imap,
		Smtp:            smtp,
	}
	settings = &result
	return
}

func GetSetting(userId *int64, key *string, defaultValue *string) (setting string) {
	appConfiguration, errAC := application.Get()
	if errAC != nil {
		log.Println(errAC.Error(), errAC)
		return *defaultValue
	}
	result, ok := appConfiguration.Get(*key)
	if ok {
		return result
	}
	return *defaultValue
}
