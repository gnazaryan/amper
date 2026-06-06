package database

import (
	databasecache "amper/cache/database"
	"amper/common/constants"
	"amper/common/structs"
	"amper/common/util"
	"amper/common/util/ampstrings"
	"container/list"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
)

const featchRecords = "SELECT %smain.* FROM amper.%s main %s WHERE %s %s %s LIMIT %d OFFSET %d"
const featchRecordsCount = "SELECT Count(main.id) AS totalCount FROM amper.%s main %s"

func FetchRecords(userId *int64, apiName *string, start *int64, limit *int64, searchParams map[string]interface{}, searchParameters *structs.Search, foreginKey *map[string]string) (result *[]map[string]interface{}, resultTotalCount int64, err error) {
	operator, okO := searchParams["operator"]
	totalCount, okTC := searchParams["totalCount"]
	in, okIn := searchParams["in"]
	operatorAnd := true
	if okO && strings.ToLower(operator.(string)) == "or" {
		operatorAnd = false
	}
	var whereClouses = []string{}
	if searchParameters != nil {
		whereClouses = *searchParameters.WhereClouses
	}
	idWhereClouse := fmt.Sprintf("main.id > %d", -1)
	if okIn {
		inMap := in.(map[string]interface{})
		if len(inMap) > 0 {
			for key, value := range inMap {
				whereClouse := fmt.Sprintf("IN (%s)", value)
				operator := ""
				if len(whereClouses) > 0 {
					operator = util.IfElse(operatorAnd, "AND", "OR").(string)
				}
				whereClouses = append(whereClouses, fmt.Sprintf("%s %s %s", operator, key, whereClouse))
			}
		}
	}
	orderBy := ""
	joins := make([]string, 0)
	joinSelects := make([]string, 0)
	if foreginKey != nil {
		for key, value := range *foreginKey {
			variable := ampstrings.RandStringBytes(3)
			join := "LEFT JOIN " + value + " " + variable + " on main." + key + " = " + variable + ".identifier_sys"
			joins = append(joins, join)
			joinSelects = append(joinSelects, "IFNULL("+variable+".name_sys, '') AS "+key+"_name_sys")

			if searchParameters.SortField != nil && searchParameters.SortDir != nil && *searchParameters.SortField == key {
				orderBy = "ORDER BY " + variable + ".name_sys" + " " + *searchParameters.SortDir
			}
		}
	}

	var pool *sql.DB = databasecache.Pool()
	whereClousesResult := ""
	if len(whereClouses) > 0 {
		whereClousesResult = " (" + strings.Join(whereClouses, " ") + ")"
	}
	joinSelectsResult := ""
	if len(joinSelects) > 0 {
		joinSelectsResult = strings.Join(joinSelects, ", ") + ", "
	}

	if searchParameters != nil && searchParameters.SortField != nil && searchParameters.SortDir != nil && len(orderBy) < 1 {
		orderBy = "ORDER BY main." + *searchParameters.SortField + " " + *searchParameters.SortDir
	}

	query := fmt.Sprintf(featchRecords, joinSelectsResult, *apiName, strings.Join(joins, " "), idWhereClouse, util.IfElse(len(whereClousesResult) > 0, "AND "+whereClousesResult, whereClousesResult), orderBy, *limit, *start)
	rows, errF := pool.Query(query)
	if errF == nil {
		cols, _ := rows.Columns()
		records := list.New()
		index := 0
		for rows.Next() {
			record := structs.Record{}
			recordData := make(map[string]interface{})
			columns := make([]string, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i, _ := range columns {
				columnPointers[i] = &columns[i]
			}

			rows.Scan(columnPointers...)

			for i, colName := range cols {
				recordData[colName] = columns[i]
			}
			id, ok := recordData[constants.GetDefaultFields().ID.ApiName]
			if ok {
				record.ID, _ = strconv.ParseInt(id.(string), 10, 64)
				record.Record = recordData
			}
			records.PushBack(record)
			index++
		}
		rows.Close()
		result = ToArray(records)

		if okTC && reflect.TypeOf(totalCount).String() == "bool" && totalCount == true {
			totalCountQuery := fmt.Sprintf(featchRecordsCount, *apiName, util.IfElse(len(whereClousesResult) > 0, "WHERE "+whereClousesResult, whereClousesResult))
			row := pool.QueryRow(totalCountQuery)
			if row.Err() == nil {
				row.Scan(&resultTotalCount)
			} else {
				log.Println(row.Err(), row.Err())
				//skip reporting the totalc count error
			}
		}
	} else {
		log.Println(errF.Error(), errF)
		err = fmt.Errorf("unable to fetch records due to errors in query parameters")
	}
	return result, resultTotalCount, err
}

func ToArray(input *list.List) (result *[]map[string]interface{}) {
	temp := make([]map[string]interface{}, input.Len())
	index := 0
	for e := input.Front(); e != nil; e = e.Next() {
		temp[index] = e.Value.(structs.Record).Record
		index++
	}
	result = &temp
	return
}

const insertRecord string = "INSERT INTO amper.%s (id, %s) VALUES (null, %s);"

func AddRecord(userId *int64, apiName *string, payload *map[string]string) (result *structs.Record, err error) {
	var pool *sql.DB = databasecache.Pool()
	var keys []string
	var values []string
	for key, value := range *payload {
		keys = append(keys, key)
		values = append(values, "'"+value+"'")
	}
	query := fmt.Sprintf(insertRecord, *apiName, strings.Join(keys, ","), strings.Join(values, ","))
	res, errDb := pool.Exec(query)
	if errDb != nil {
		log.Println(errDb.Error(), errDb)
		err = fmt.Errorf("unable to insert a record into the database due to sql error")
		return
	} else if count, _ := res.RowsAffected(); count < 1 {
		err = fmt.Errorf("unable to insert a record into the database")
		return
	}
	id, _ := res.LastInsertId()
	newRecord, err := FetchRecord(userId, apiName, &id, nil)
	result = &newRecord
	return result, err
}

const selectRecord = "Select * from amper.%s WHERE %s"

func FetchRecord(userId *int64, apiName *string, id *int64, identifier *string) (result structs.Record, err error) {
	var pool *sql.DB = databasecache.Pool()
	whereClouse := ""
	if identifier != nil {
		whereClouse = constants.GetDefaultFields().IDENTIFIER.ApiName + "='" + *identifier + "'"
	} else if id != nil {
		whereClouse = constants.GetDefaultFields().ID.ApiName + "='" + strconv.FormatInt(*id, 10) + "'"
	} else {
		err = fmt.Errorf("unable to fetch a record due to missing identifying information")
		return
	}
	rows, err := pool.Query(fmt.Sprintf(selectRecord, *apiName, whereClouse))
	if err == nil {
		cols, _ := rows.Columns()
		record := make(map[string]interface{})

		if rows.Next() {
			columns := make([]string, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i, _ := range columns {
				columnPointers[i] = &columns[i]
			}

			rows.Scan(columnPointers...)

			for i, colName := range cols {
				record[colName] = columns[i]
			}
			id, ok := record[constants.GetDefaultFields().ID.ApiName]
			if ok {
				result.ID, _ = strconv.ParseInt(id.(string), 10, 64)
				result.Record = record
			}
		}
		rows.Close()
	} else {
		log.Print(err.Error(), err)
		err = fmt.Errorf("unable to fetch a record with object api name %s and id %d", *apiName, id)
	}
	return
}

const deleteRecord = "DELETE FROM amper.%s WHERE %s"

func RemoveRecord(userId *int64, apiName *string, id *int64, identifier *string) (result *structs.Record, err error) {
	record, errR := FetchRecord(userId, apiName, id, identifier)
	if errR == nil && record.ID > 0 {
		whereClouse := ""
		if identifier != nil {
			whereClouse = constants.GetDefaultFields().IDENTIFIER.ApiName + "='" + *identifier + "'"
		} else if id != nil {
			whereClouse = constants.GetDefaultFields().ID.ApiName + "='" + strconv.FormatInt(*id, 10) + "'"
		} else {
			err = fmt.Errorf("unable to remove a record due to missing identifying information")
			return
		}
		var pool *sql.DB = databasecache.Pool()
		rows, errD := pool.Exec(fmt.Sprintf(deleteRecord, *apiName, whereClouse))
		if errD != nil {
			var errorCode uint16 = 0
			if reflect.TypeOf(errD).String() == "*mysql.MySQLError" {
				var interfaceError interface{} = errD
				mysqlError := interfaceError.(*mysql.MySQLError)
				errorCode = mysqlError.Number
			}
			log.Println(errD.Error(), errD)
			err = fmt.Errorf("[%d] unable to remove a record with id %d and identifier %s", errorCode, id, *identifier)
			return
		} else if count, _ := rows.RowsAffected(); count > 0 {
			result = &record
			return
		}
	} else {
		if errR != nil {
			log.Println(errR.Error(), errR)
		}
		err = fmt.Errorf("unable to remove a record with id %d and identifier %s, no record exist with the provided identifier", id, util.IfElse(identifier != nil, identifier, ampstrings.EmptyIfNil(identifier)))
	}
	return result, err
}

func RemoveRecordById(userId *int64, id int64) (result *structs.Record, err error) {
	return result, nil
}

const updateRecord = "UPDATE amper.%s SET %s WHERE %s"

func UpdateRecord(userId *int64, apiName *string, payload *map[string]string) (result bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	var items []string
	for key, value := range *payload {
		if constants.GetDefaultFields().ID.ApiName != key {
			items = append(items, " "+key+"='"+value+"'")
		}
	}
	id, idOk := (*payload)[constants.GetDefaultFields().ID.ApiName]
	identifier, identifierOk := (*payload)[constants.GetDefaultFields().IDENTIFIER.ApiName]
	whereClous := ""
	if idOk {
		whereClous = constants.GetDefaultFields().ID.ApiName + "='" + id + "'"
	} else if identifierOk {
		whereClous = constants.GetDefaultFields().IDENTIFIER.ApiName + "='" + identifier + "'"
	} else {
		err = fmt.Errorf("unable to update a record in the database, identifier is not supplied")
		return
	}
	query := fmt.Sprintf(updateRecord, *apiName, strings.Join(items, ","), whereClous)
	res, errDb := pool.Exec(query)
	if errDb != nil {
		log.Println(errDb.Error(), errDb)
		err = fmt.Errorf("unable to update a record in the database due to sql error")
		return
	} else if count, _ := res.RowsAffected(); count < 1 {
		err = fmt.Errorf("unable to update a record in the database")
		return
	}
	return true, nil
}

const insertRecords = "INSERT INTO amper.%s (%s) VALUES %s"
const selectLastInserted = "SELECT * FROM amper.%s WHERE id > %d AND id < %d"

func AddRecords(userId *int64, apiName *string, payloadPartition *list.List) (resultSuccess list.List, err error) {
	var values strings.Builder
	keys := make([]string, 0)
	keysLoaded := false
	for payload := payloadPartition.Front(); payload != nil; payload = payload.Next() {
		payloadData := payload.Value.(map[string]string)
		if !keysLoaded {
			for key, _ := range payloadData {
				keys = append(keys, key)
			}
			keysLoaded = true
		}
		firstRun := true
		if values.Len() > 1 {
			values.WriteString(", ")
		}
		values.WriteString("(")
		for _, key := range keys {
			if !firstRun {
				values.WriteString(", ")
			} else {
				firstRun = false
			}
			values.WriteString("'")
			values.WriteString(payloadData[key])
			values.WriteString("'")
		}
		values.WriteString(")")
	}
	query := fmt.Sprintf(insertRecords, *apiName, strings.Join(keys, ","), values.String())
	var pool *sql.DB = databasecache.Pool()
	res, errI := pool.Exec(query)
	if errI != nil {
		log.Print(errI.Error(), errI)
		err = fmt.Errorf("there was a databse execution error and the record addition was failed")
	} else {
		lastInsertId, _ := res.LastInsertId()
		rowsAffected, _ := res.RowsAffected()
		selectQuery := fmt.Sprintf(selectLastInserted, *apiName, lastInsertId-1, lastInsertId+rowsAffected)
		rows, errLI := pool.Query(selectQuery)
		if errLI == nil {
			cols, _ := rows.Columns()
			data := structs.Record{}
			for rows.Next() {
				record := make(map[string]interface{})
				columns := make([]string, len(cols))
				columnPointers := make([]interface{}, len(cols))
				for i, _ := range columns {
					columnPointers[i] = &columns[i]
				}

				rows.Scan(columnPointers...)

				for i, colName := range cols {
					record[colName] = columns[i]
				}
				id, ok := record[constants.GetDefaultFields().ID.ApiName]
				if ok {
					data.ID, _ = strconv.ParseInt(id.(string), 10, 64)
					data.Record = record
				}
				resultSuccess.PushBack(data)
			}
			rows.Close()
		} else {
			log.Print(errI.Error(), errI)
			err = fmt.Errorf("there was a databse execution error and the record addition was failed")
		}
	}
	return resultSuccess, err
}

const updateRecordsIdentifiers = "Update amper.%s SET identifier_sys = TO_BASE64(CONCAT(SUBSTRING_INDEX(identifier_sys, '|' , 3), '|', id)) WHERE id IN (%s)"

func ComputeRecordsIdentifier(userId *int64, apiName *string, partionIds *[]int64) (result bool, err error) {
	if partionIds == nil {
		return true, nil
	}
	join := ampstrings.JoinInt64(partionIds, ", ")
	query := fmt.Sprintf(updateRecordsIdentifiers, *apiName, *join)
	var pool *sql.DB = databasecache.Pool()
	res, errC := pool.Exec(query)
	if errC != nil {
		log.Print(errC.Error(), errC)
		err = fmt.Errorf("there was a databse execution error and the record addition was failed")
		result = false
	} else if count, _ := res.RowsAffected(); count != int64(len(*partionIds)) {
		err = fmt.Errorf("there was a databse execution error and number of affected rows is not same as requested for object: %s, ids: %s", *apiName, *join)
		result = false
	} else {
		result = true
	}
	return result, err
}

const updateRecords = "UPDATE amper.%s SET %s WHERE id IN (%s)"

func UpdateRecords(userId *int64, apiName *string, partition *list.List) (result bool, err error) {
	if partition == nil || partition.Len() < 1 {
		return true, nil
	}
	wheneClouses := make(map[string]*strings.Builder)
	ids := new(strings.Builder)
	for payload := partition.Front(); payload != nil; payload = payload.Next() {
		payloadData := payload.Value.(*map[string]string)
		id, idOk := (*payloadData)["id"]
		if idOk && len(id) > 0 {
			if ids.Len() > 0 {
				ids.WriteString(",")
			}
			ids.WriteString(id)
			for key, value := range *payloadData {
				if key == "id" {
					continue
				}
				wheneClouse := wheneClouses[key]
				if wheneClouse == nil {
					wheneClouse = new(strings.Builder)
					wheneClouse.WriteString(key)
					wheneClouse.WriteString(" = CASE")
					wheneClouses[key] = wheneClouse
				}
				wheneClouse.WriteString(" WHEN id = ")
				wheneClouse.WriteString(id)
				wheneClouse.WriteString(" THEN '")
				wheneClouse.WriteString(value)
				wheneClouse.WriteString("'")
			}
		}
	}
	whenClousesResult := new(strings.Builder)
	for key, wheneClouse := range wheneClouses {
		wheneClouse.WriteString(" ELSE ")
		wheneClouse.WriteString(key)
		wheneClouse.WriteString(" END")
		if whenClousesResult.Len() > 0 {
			whenClousesResult.WriteString(", ")
		}
		whenClousesResult.WriteString(wheneClouse.String())
	}
	query := fmt.Sprintf(updateRecords, *apiName, whenClousesResult.String(), ids.String())
	var pool *sql.DB = databasecache.Pool()
	rows, errD := pool.Exec(query)
	if errD != nil || rows == nil {
		util.Loggify(errD)
		err = fmt.Errorf("unable to update records with supplied ids: %s for object %s", ids.String(), *apiName)
	} else if count, _ := rows.RowsAffected(); count > 0 {
		result = true
	}
	return result, err
}

const removeRecords = "DELETE FROM amper.%s WHERE id IN (%s)"

func RemoveRecords(userId *int64, apiName *string, ids *list.List) (result bool, err error) {
	if ids == nil || ids.Len() < 1 {
		return true, nil
	}
	idsComma := ampstrings.JoinListInt64(ids, ", ")
	query := fmt.Sprintf(removeRecords, *apiName, *idsComma)
	var pool *sql.DB = databasecache.Pool()
	rows, errD := pool.Exec(query)
	if errD != nil || rows == nil {
		var errorCode uint16 = 0
		if reflect.TypeOf(errD).String() == "*mysql.MySQLError" {
			var interfaceError interface{} = errD
			mysqlError := interfaceError.(*mysql.MySQLError)
			errorCode = mysqlError.Number
		}

		util.Loggify(errD)
		err = fmt.Errorf("[%d] unable to remove records with supplied ids: %s for object %s", errorCode, *idsComma, *apiName)
	} else if count, _ := rows.RowsAffected(); count > 0 {
		result = true
	}
	return result, err
}
