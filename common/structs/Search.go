package structs

import (
	"amper/common/util"
	"amper/common/util/arrays"
	"amper/common/util/jsons"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type Search struct {
	WhereClouses *[]string
	Operator     *string
	Errors       *[]string

	SortField *string
	SortDir   *string
}

func (r *Search) Parse(search *string) (result bool) {
	whereClouses := make([]string, 0)
	errors := make([]string, 0)
	operator := "AND"
	var searchParams map[string]interface{}
	if search != nil && len(*search) > 0 {
		var errJ error
		searchParams, errJ = jsons.GetJsonObject(search)
		if errJ != nil {
			log.Print(errJ.Error(), errJ)
			errors = append(errors, "unable to fetch records, the supplied search data is in a wrong format")
		}
	}

	if len(searchParams) > 0 {
		sortField, okSF := searchParams["sortField"]
		if okSF && reflect.TypeOf(sortField).String() == "string" {
			r.SortField = util.PointerString(sortField.(string))
		}
		sortDir, okSF := searchParams["sortDir"]
		if okSF && reflect.TypeOf(sortDir).String() == "string" {
			r.SortDir = util.PointerString(sortDir.(string))
		}
		operator, okO := searchParams["operator"]
		if okO && strings.ToLower(operator.(string)) == "or" {
			operator = "OR"
		} else {
			operator = "AND"
		}
		term, okT := searchParams["term"]
		if okT {
			termMap, okTM := term.([]interface{})
			if okTM && len(termMap) > 0 {
				for _, filterMap := range termMap {
					filter, okF := filterMap.(map[string]interface{})
					if okF {
						operatorCondition, okO := filter["operator"].(string)
						apiName, okA := filter["apiName"].(string)
						if okO && okA && len(operatorCondition) > 0 {
							switch operatorCondition {
							case "greaterThen":
								value, errV := util.I2Num(filter["value"])
								if errV == nil {
									whereClouse := fmt.Sprintf("> %d", value)
									whereClouses = append(whereClouses, fmt.Sprintf("%s %s %s",
										util.IfElse(len(whereClouses) > 0, operator, ""), "main."+apiName, whereClouse))
								} else {
									errors = append(errors, "the supplied operator 'greaterThen' has no properly formatted number value, the value parameter is required")
								}
							case "lessThen":
								value, errV := util.I2Num(filter["value"])
								if errV == nil {
									whereClouse := fmt.Sprintf("< %d", value)
									whereClouses = append(whereClouses, fmt.Sprintf("%s %s %s",
										util.IfElse(len(whereClouses) > 0, operator, ""), "main."+apiName, whereClouse))
								} else {
									errors = append(errors, "the supplied operator 'lessThen' has no properly formatted number value, the value parameter is required")
								}
							case "equalsNumber":
								value, errV := util.I2Num(filter["value"])
								if errV == nil {
									whereClouses = append(whereClouses, fmt.Sprintf("%s %s=%d", util.IfElse(len(whereClouses) > 0, operator, ""), "main."+apiName, value))
								} else {
									errors = append(errors, "the supplied operator 'equalsNumber' has no properly formatted number value, the value parameter is required")
								}
							case "notEqualsNumber":
								value, errV := util.I2Num(filter["value"])
								if errV == nil {
									whereClouses = append(whereClouses, fmt.Sprintf("%s %s<>%d", util.IfElse(len(whereClouses) > 0, operator, ""), "main."+apiName, value))
								} else {
									errors = append(errors, "the supplied operator 'notEqualsNumber' has no properly formatted number value, the value parameter is required")
								}
							case "isNotEmptyNumber":
								whereClouses = append(whereClouses, fmt.Sprintf("%s %s IS NOT NULL", util.IfElse(len(whereClouses) > 0, operator, ""), "main."+apiName))
							case "isEmptyNumber":
								whereClouses = append(whereClouses, fmt.Sprintf("%s %s IS NULL", util.IfElse(len(whereClouses) > 0, operator, ""), "main."+apiName))
							case "hasAnyOf":
								valuesInterface, okInt := filter["value"].([]interface{})
								if okInt {
									values, okV := arrays.InterfaceToString(&valuesInterface)
									if okV {
										if len(values) > 0 {
											whereClouse := fmt.Sprintf("IN ('%s')", strings.Join(values, "', '"))
											whereClouses = append(whereClouses, fmt.Sprintf("%s %s %s",
												util.IfElse(len(whereClouses) > 0, operator, ""), "main."+apiName, whereClouse))
										}
									} else {
										errors = append(errors, "the supplied operator 'hasAnyOf' has no properly formatted array of values, the value parameter is required")
									}
								} else {
									errors = append(errors, "the supplied operator 'hasAnyOf' has no properly formatted array of values, the value parameter is required")
								}
							case "hasNoneOf":
								valuesInterface, okInt := filter["value"].([]interface{})
								if okInt {
									values, okV := arrays.InterfaceToString(&valuesInterface)
									if okV {
										if len(values) > 0 {
											whereClouse := fmt.Sprintf("NOT IN ('%s')", strings.Join(values, "', '"))
											whereClouses = append(whereClouses, fmt.Sprintf("%s %s %s",
												util.IfElse(len(whereClouses) > 0, operator, ""), "main."+apiName, whereClouse))
										}
									} else {
										errors = append(errors, "the supplied operator 'hasNoneOf' has no properly formatted array of values, the value parameter is required")
									}
								} else {
									errors = append(errors, "the supplied operator 'hasNoneOf' has no properly formatted array of values, the value parameter is required")
								}
							case "contains":
								value, okV := filter["value"].(string)
								if okV {
									if len(value) > 0 {
										whereClouse := fmt.Sprintf("REGEXP '%s'", value)
										whereClouses = append(whereClouses, fmt.Sprintf("%s %s %s", util.IfElse(len(whereClouses) > 0, operator, ""), "main."+apiName, whereClouse))
									}
								} else {
									errors = append(errors, "the supplied operator 'contains' has no properly formatted string value, the value parameter is required")
								}
							case "equals":
								value, okV := filter["value"].(string)
								if okV {
									if len(value) > 0 {
										whereClouses = append(whereClouses, fmt.Sprintf("%s %s = '%s'", util.IfElse(len(whereClouses) > 0, operator, ""), "main."+apiName, value))
									}
								} else {
									errors = append(errors, "the supplied operator 'equals' has no properly formatted string value, the value parameter is required")
								}
							case "startsWith":
								value, okV := filter["value"].(string)
								if okV {
									if len(value) > 0 {
										whereClouses = append(whereClouses, fmt.Sprintf("%s %s LIKE '%%%s'", util.IfElse(len(whereClouses) > 0, operator, ""), "main."+apiName, value))
									}
								} else {
									errors = append(errors, "the supplied operator 'startsWith' has no properly formatted string value, the value parameter is required")
								}
							case "endsWith":
								value, okV := filter["value"].(string)
								if okV {
									if len(value) > 0 {
										whereClouses = append(whereClouses, fmt.Sprintf("%s %s LIKE '%s%%'", util.IfElse(len(whereClouses) > 0, operator, ""), "main."+apiName, value))
									}
								} else {
									errors = append(errors, "the supplied operator 'endsWith' has no properly formatted string value, the value parameter is required")
								}
							case "isEmpty":
								column := "main." + apiName
								whereClouses = append(whereClouses, fmt.Sprintf("%s (%s IS NULL OR %s = '')", util.IfElse(len(whereClouses) > 0, operator, ""), column, column))
							case "isNotEmpty":
								column := "main." + apiName
								whereClouses = append(whereClouses, fmt.Sprintf("%s (%s IS NOT NULL OR %s <> '')", util.IfElse(len(whereClouses) > 0, operator, ""), column, column))
							case "is":
								value, okV := filter["value"].(bool)
								if okV {
									sqlValue := int8(1)
									if !value {
										sqlValue = int8(0)
									}
									whereClouses = append(whereClouses, fmt.Sprintf("%s %s=%d", util.IfElse(len(whereClouses) > 0, operator, ""), "main."+apiName, sqlValue))
								} else {
									errors = append(errors, "the supplied operator 'is' has no properly formatted string value, the value parameter is required")
								}
							default:
								errors = append(errors, "the supplied operator in term filter has no value, the operator parameter is required")
							}
						} else {
							errors = append(errors, "the supplied operator '%s' is not supported red the manual and try a different one")
						}
					} else {
						errors = append(errors, "unable to parse search parameter, the enclosed term has to be array of json objects")
					}
				}
			}
		}
	}
	result = len(errors) == 0
	r.WhereClouses = &whereClouses
	r.Operator = &operator
	r.Errors = &errors
	return result
}
