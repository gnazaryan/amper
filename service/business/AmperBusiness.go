package business

import (
	"amper/cache/business"
	"amper/common/structs"
	"amper/common/util"
	"amper/common/util/arrays"
	"amper/common/util/datetime"
	"amper/data/database"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func init() {
	EditLastUpdateTime()
}

func EditLastUpdateTime() {
	time.AfterFunc(5*time.Minute, func() {
		database.EditLastUpdateTime(nil, business.AmperId(), datetime.GetDateTimeFormatted())
		EditLastUpdateTime()
	})
}

func GetInstances(userID *int64, Type *string) ([]structs.Amper, error) {
	instances, err := database.GetInstances(userID, Type, false)
	for i := 0; i < len(instances); i++ {
		if instances[i].StateUpdateDate != nil {
			lastUpdateDate, errLUD := datetime.ParseDateTime(instances[i].StateUpdateDate)
			now, _ := datetime.ParseDateTime(util.PointerString(datetime.GetDateTimeFormatted()))
			diff := now.Sub(lastUpdateDate)
			if errLUD != nil || diff.Minutes() > 6 {
				instances[i].State = util.PointerInt(0)
			}
		} else {
			instances[i].State = util.PointerInt(0)
		}
	}
	return instances, err
}

func RemoveInstance(userID *int64, amper structs.Amper) (bool, error) {
	success, err := database.RemoveInstance(userID, amper)
	return success, err
}

func CreateInstance(userID *int64, sessionId *string, amper structs.Amper) (success bool, err error) {
	confirmedAmper, errA := FetchInstanceInfo(userID, sessionId, amper)
	if errA == nil {
		success, err = database.CreateInstance(userID, *confirmedAmper)
		if !success || err != nil {
			util.Loggify(err)
			err = fmt.Errorf("not able to create a database record for amper instance, communicate with support")
		}
	} else {
		util.Loggify(errA)
		err = fmt.Errorf("not able to verify the amper instance before creating, make sure the adress and port are correct")
	}

	return success, err
}

func EditInstance(userID *int64, sessionId *string, amper structs.Amper) (bool, error) {
	success, err := database.EditInstance(userID, amper)
	if success && err == nil {
		instanceTypes := []string{"amperInstance", "amperDatastoreInstance"}
		successFC, errFC := FederatedCall(userID, sessionId, map[string]string{
			"amperDatastoreInstance": "amper-datastore/invalidateCache",
			"amperInstance":          "amper/invalidateCache",
		}, map[string]interface{}{
			"name": "amper",
		}, &instanceTypes)
		if !successFC || errFC != nil {
			util.Loggify(errFC)
			err = fmt.Errorf("not able to reset the remote cache for amper, contact the support")
			success = false
		}
	}
	return success, err
}

func FetchInstanceInfo(userID *int64, sessionId *string, amper structs.Amper) (result *structs.Amper, err error) {
	body, err := json.Marshal(amper)
	url := util.IfElse(*amper.Type == "amperDatastoreInstance", "http://%s:%s/amper-datastore/status", "http://%s:%s/amper/status")
	r, errNR := http.NewRequest("POST", fmt.Sprintf(url.(string), *amper.Address, *amper.Port), bytes.NewBuffer(body))
	if errNR == nil {
		r.Header.Add("userId", strconv.FormatInt(*userID, 10))
		r.Header.Add("sessionId", *sessionId)

		client := &http.Client{}
		res, errDo := client.Do(r)
		if errDo == nil {
			defer res.Body.Close()
			post := &structs.AmperResult{}
			errD := json.NewDecoder(res.Body).Decode(post)
			if errD == nil {
				if post.Success {
					if *post.Data.Type == *amper.Type {
						amper.Identifier = post.Data.Identifier
						amper.State = post.Data.State
						return &amper, nil
					} else {
						err = fmt.Errorf("the specified amper instance is of a different type: %s", *post.Data.Type)
					}
				} else {
					err = fmt.Errorf("the specified amper instance directory is not valid: %s", *amper.Directory)
				}
			} else {
				log.Println(errD.Error(), errD)
				err = fmt.Errorf("not able to decode the response from the amper instance, try a different address or port")
			}
		} else {
			log.Println(errDo.Error(), errDo)
			err = fmt.Errorf("not able to reach the amper instance, try a different address or port")
		}
	} else {
		log.Println(errNR.Error(), errNR)
		err = fmt.Errorf("not able to reach the amper instance, please try again with different address and port or reach the support")
	}
	return result, err
}

func FetchStatus(userID *int64, amper structs.Amper) (result *structs.Amper, err error) {
	result = &structs.Amper{}
	_, errS := os.Stat(*amper.Directory)
	if errS == nil {
		result.State = util.PointerInt(1)
	} else {
		result.State = util.PointerInt(0)
		err = fmt.Errorf("the suplied directory doesn't exist, try something different")
	}
	result.Type = util.PointerString("amperInstance")
	result.Identifier = business.AmperId()
	return result, err
}

// The function is designed to make a http request to the specified amper instance with a retry attempts
// !!!ATTENTION!!! The same copy of code is used inside the datastore repository, if this method is changed or fixed
// Consider making the changes or fixing the issues in the datastore repo too
func DedicatedCallWithRetry(userID *int64, sessionId *string, api map[string]string, parameters map[string]interface{}, instance *structs.Amper) (success bool, value *string, err error) {
	success = true
	for i := 0; i < retryAttempts; i++ {
		partialSuccess, partialValue, retry, errD := DedicatedCall(userID, sessionId, api, parameters, instance)
		util.Loggify(errD)
		if partialSuccess || !retry {
			success = partialSuccess
			value = partialValue
			break
		}
		if i+1 == retryAttempts {
			success = false
		}
	}
	return success, value, err
}

// The function is designed to make a http request to the specified amper instance
// !!!ATTENTION!!! The same copy of code is used inside the datastore repository, if this method is changed or fixed
// Consider making the changes or fixing the issues in the datastore repo too
func DedicatedCall(userID *int64, sessionId *string, api map[string]string, parameters map[string]interface{}, instance *structs.Amper) (success bool, value *string, retry bool, err error) {
	parametersString, _ := json.Marshal(parameters)
	body := []byte(parametersString)

	url := fmt.Sprintf("http://%s:%s/%s", *instance.Address, *instance.Port, api[*instance.Type])
	r, errNR := http.NewRequest("POST", url, bytes.NewBuffer(body))
	retry = true
	if errNR == nil {
		r.Header.Add("userId", strconv.FormatInt(*userID, 10))
		if sessionId != nil {
			r.Header.Add("sessionId", *sessionId)
		}

		client := &http.Client{}
		res, errDo := client.Do(r)
		if errDo == nil {
			defer res.Body.Close()
			post := &structs.ResultValue{}
			errD := json.NewDecoder(res.Body).Decode(post)
			if errD == nil {
				if post.Success {
					success = true
					value = post.Value
				} else {
					err = fmt.Errorf("node `%s` received failed internal responce with error: %s", url, post.Error)
					retry = false
				}
			} else {
				util.Loggify(errD)
				err = fmt.Errorf("node `%s` failed: not able to decode the response from the amper instance, try a different address or port", url)
			}
		} else {
			util.Loggify(errDo)
			err = fmt.Errorf("node `%s` failed: not able to reach the amper instance, try a different address or port", url)
		}
	} else {
		util.Loggify(errNR)
		err = fmt.Errorf("node `%s` failed: not able to reach the amper instance, please try again with different address and port or reach the support", url)
	}

	return success, value, retry, err
}

var retryAttempts int = 3

// The function is designed to make a federated http request across all nodes with a retry
// !!!ATTENTION!!! The same copy of code is used inside the datastore repository, if this method is changed or fixed
// Consider making the changes or fixing the issues in the datastore repo too
func FederatedCall(userID *int64, sessionId *string, api map[string]string, parameters map[string]interface{}, nodeTypes *[]string) (success bool, err error) {
	instances, errI := database.GetInstances(userID, nil, false)
	if errI != nil {
		util.Loggify(errI)
		return false, fmt.Errorf("not able to fetch instances for federated call")
	}
	success = true
	for i := 0; i < len(instances); i++ {
		if !arrays.ContainsS(*nodeTypes, *instances[i].Type) {
			continue
		}
		partialSuccess := true
		for l := 0; l < retryAttempts; l++ {
			dedicatedSuccess, _, retry, errI := DedicatedCall(userID, sessionId, api, parameters, &instances[i])
			util.Loggify(errI)
			if dedicatedSuccess || !retry {
				partialSuccess = dedicatedSuccess
				break
			}
			if l+1 == retryAttempts {
				partialSuccess = false
			}
		}
		if !partialSuccess {
			success = false
		}
	}

	return success, err
}

func InvalidateCache(userID *int64, Name *string, userIdDelete *int64, chatChannelDelete *int64) (success bool, err error) {
	switch *Name {
	case "amper":
		business.InvalidateAmperCache()
		success = true
	case "user":
		business.InvalidateUserCache(userIdDelete)
		success = true
	case "chatChannel":
		business.InvalidateChatChannel(chatChannelDelete)
		success = true
	}
	return success, err
}

func GetAvailableAmperDatastore(userID *int64) (*structs.Amper, error) {
	instances, errI := GetInstances(userID, util.PointerString("amperDatastoreInstance"))
	if errI != nil {
		util.Loggify(errI)
		return nil, fmt.Errorf("not able to fetch datastore instances from database")
	}
	var mostAvailableSpace int64 = 0
	var mostAvalableInstance *structs.Amper
	for _, instance := range instances {
		if *instance.Limit-*instance.Usage > mostAvailableSpace {
			mostAvalableInstance = &instance
			mostAvailableSpace = *instance.Limit - *instance.Usage
		}
	}
	return mostAvalableInstance, nil
}
