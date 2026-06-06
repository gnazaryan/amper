package structs

import (
	"amper/common/util"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
)

type Record struct {
	ID     int64                  `json:"id"`
	Record map[string]interface{} `json:"record"`
}

type RecordsResult struct {
	Result
	Data       *[]map[string]interface{} `json:"data"`
	ErrorData  *[]map[string]interface{} `json:"errordata"`
	Metadata   *ObjectMetadata           `json:"metadata"`
	TotalCount int64                     `json:"totalCount"`
}

type RecordResult struct {
	Result
	Data *Record `json:"data"`
}

type RecordInfo struct {
	Identifier   *string
	Id           *int64
	ObjectTypeId *int64
	ObjectId     *int64
	AmperId      *int64
}

func (r Record) geStringtValue(apiName *string) string {
	var result string
	if r.Record[*apiName] != nil {
		result = r.Record[*apiName].(string)
	}
	return result
}

func (r Record) getIntValue(apiName *string) int64 {
	var result int64
	if r.Record[*apiName] != nil {
		result = r.Record[*apiName].(int64)
	}
	return result
}

func (r Record) getBoolValue(apiName *string) bool {
	var result bool
	if r.Record[*apiName] != nil {
		result = r.Record[*apiName].(bool)
	}
	return result
}

func (r *RecordInfo) Parse() (result bool, err error) {
	decoded, errB := base64.StdEncoding.DecodeString(*r.Identifier)
	if errB == nil {
		identiferParts := strings.Split(string(decoded), "|")
		if len(identiferParts) == 4 {
			r.Id = util.ParseInt64(&identiferParts[3])
			r.ObjectTypeId = util.ParseInt64(&identiferParts[2])
			r.ObjectId = util.ParseInt64(&identiferParts[1])
			r.AmperId = util.ParseInt64(&identiferParts[0])
			result = true
		} else {
			err = fmt.Errorf("supplied identifer %s is wrongly formatted", *r.Identifier)
		}
	} else {
		log.Print(errB.Error(), errB)
		err = fmt.Errorf("supplied identifer %s is not valid", *r.Identifier)
	}
	return result, err
}
