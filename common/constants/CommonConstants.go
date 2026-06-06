package constants

import (
	"math"
	"reflect"
)

var API_SUFFIX = "_amp"

var TIME_FORMAT = "2006-01-02 15:04:05"

type FieldType struct {
	Name            string `json:"name"`
	Size            int64  `json:"size"`
	CreateStatement string `json:"createStatement"`
}

type DataTypes struct {
	ID         FieldType
	IDENTIFIER FieldType
	OBJECTTYPE FieldType
	TEXT       FieldType
	NUMBER     FieldType
	BOOLEAN    FieldType
	DATE       FieldType
	DATETIME   FieldType
	REFERENCE  FieldType
}

var GetDataTypes = DataTypes{
	ID: FieldType{
		Name:            "INTEGER",
		Size:            math.MaxInt32,
		CreateStatement: "`%s` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT",
	},
	IDENTIFIER: FieldType{
		Name:            "TEXT",
		Size:            128,
		CreateStatement: "`%s` varchar(%d) %s",
	},
	OBJECTTYPE: FieldType{
		Name:            "TEXT",
		Size:            128,
		CreateStatement: "`%s` varchar(%d) %s",
	},
	TEXT: FieldType{
		Name:            "TEXT",
		Size:            512,
		CreateStatement: "`%s` varchar(%d) %s",
	},
	NUMBER: FieldType{
		Name:            "NUMBER",
		Size:            math.MaxInt32,
		CreateStatement: "`%s` BIGINT(20) %s",
	},
	BOOLEAN: FieldType{
		Name:            "BOOLEAN",
		Size:            4,
		CreateStatement: "`%s` tinyint(4) %s",
	},
	DATE: FieldType{
		Name:            "DATE",
		Size:            math.MaxInt32,
		CreateStatement: "`%s` DATE %s",
	},
	DATETIME: FieldType{
		Name:            "DATETIME",
		Size:            math.MaxInt32,
		CreateStatement: "`%s` DATETIME %s",
	},
	REFERENCE: FieldType{
		Name:            "REFERENCE",
		Size:            128,
		CreateStatement: "`%s` varchar(128) %s",
	},
}

func GetFieldType(dataTypeString *string) *FieldType {
	var dataTypes = [9]FieldType{GetDataTypes.ID, GetDataTypes.IDENTIFIER, GetDataTypes.OBJECTTYPE, GetDataTypes.TEXT, GetDataTypes.NUMBER, GetDataTypes.BOOLEAN, GetDataTypes.DATE, GetDataTypes.DATETIME, GetDataTypes.REFERENCE}
	for _, dataType := range dataTypes {
		if dataType.Name == *dataTypeString {
			return &dataType
		}
	}
	return nil
}

type Field struct {
	ApiName  string    `json:"apiName"`
	Label    string    `json:"label"`
	Type     FieldType `json:"type"`
	Status   int       `json:"status"`
	Required int       `json:"required"`
}

type DefaultFields struct {
	ID         Field
	IDENTIFIER Field
	OBJECTTYPE Field
	NAME       Field
	STATUS     Field
}

func GetDefaultFields() (result DefaultFields) {
	return DefaultFields{
		ID: Field{
			ApiName:  "id",
			Label:    "Id",
			Type:     GetDataTypes.ID,
			Status:   1,
			Required: 1,
		},
		IDENTIFIER: Field{
			ApiName:  "identifier_sys",
			Label:    "Identifier",
			Type:     GetDataTypes.IDENTIFIER,
			Status:   1,
			Required: 1,
		},
		OBJECTTYPE: Field{
			ApiName:  "objectType_sys",
			Label:    "Object type",
			Type:     GetDataTypes.OBJECTTYPE,
			Status:   1,
			Required: 1,
		},
		NAME: Field{
			ApiName:  "name_sys",
			Label:    "Name",
			Type:     GetDataTypes.TEXT,
			Status:   1,
			Required: 1,
		},
		STATUS: Field{
			ApiName:  "status_sys",
			Label:    "Status",
			Type:     GetDataTypes.BOOLEAN,
			Status:   1,
			Required: 1,
		},
	}
}

func IsSystemField(apiName *string) bool {
	v := reflect.ValueOf(GetDefaultFields())
	for i := 0; i < v.NumField(); i++ {
		currentApiName := (v.Field(i).Interface().(Field)).ApiName
		if currentApiName == *apiName {
			return true
		}
	}
	return false
}

type BaseObjectType struct {
	KEY   string
	LABEL string
}

func GetBaseObjectType() (result BaseObjectType) {
	return BaseObjectType{
		KEY:   "base_sys",
		LABEL: "Base",
	}
}
