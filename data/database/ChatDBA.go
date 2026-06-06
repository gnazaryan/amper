package database

import (
	databasecache "amper/cache/database"
	"amper/common/structs"
	"amper/common/util"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
)

const insertChatChannelGroup = "INSERT INTO amper.chat_channel_group_sys VALUES (null, '%s')"

func CreateChatGroup(userID *int64, Name *string) (success bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(insertChatChannelGroup, *Name)
	res, errW := pool.Exec(query)
	if errW != nil {
		log.Print(errW.Error(), errW)
		return false, fmt.Errorf("unable to insert an channel group into the database with specified parameters: name - %s", *Name)
	} else if count, _ := res.RowsAffected(); count < 1 {
		return false, fmt.Errorf("unable to insert an amper instance into the database with specified parameters: name - %s", *Name)
	}
	return true, nil
}

const getChatChannelGroupQuery = "SELECT * FROM amper.chat_channel_group_sys"

func FetchChatChannelGroups(userID *int64) (result []structs.ChatChannelGroup, err error) {
	var pool *sql.DB = databasecache.Pool()
	rows, errQ := pool.Query(getChatChannelGroupQuery)
	if errQ == nil {
		for rows.Next() {
			var channelGroup structs.ChatChannelGroup
			rows.Scan(&channelGroup.Id, &channelGroup.Name)
			result = append(result, channelGroup)
		}
	} else {
		err = errors.New("unable to run query against database to get chat channel groups")
		log.Print(errQ.Error(), errQ)
	}
	rows.Close()
	return result, err
}

const insertChatChannel = "INSERT INTO amper.chat_channel_sys VALUES (null, '%s', '%d', '%d', '', '')"

func CreateChatChannel(userID *int64, Name *string, AmperId *int64, GroupId *int64) (id int64, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(insertChatChannel, *Name, *GroupId, *AmperId)
	res, errW := pool.Exec(query)
	if errW != nil {
		log.Print(errW.Error(), errW)
		return -1, fmt.Errorf("unable to insert an channel into the database with specified parameters: name - %s", *Name)
	} else if count, _ := res.RowsAffected(); count < 1 {
		return -1, fmt.Errorf("row cound is 0, unable to insert an channel into the database with specified parameters: name - %s", *Name)
	}
	lastInsertId, errLII := res.LastInsertId()
	if errLII != nil {
		return -1, fmt.Errorf("no last insert id, unable to insert an channel into the database with specified parameters: name - %s", *Name)
	}
	return lastInsertId, nil
}

const getChatChannelsQuery = "SELECT cc.id as id, cc.name as name, cc.amper_id as amperId, a.name as amperName, cc.group_id as groupId, ccg.name as groupName FROM amper.chat_channel_sys as cc inner join amper.amper_sys as a on a.id = cc.amper_id inner join amper.chat_channel_group_sys as ccg on ccg.id = cc.group_id"

func FetchChatChannels(userID *int64, GroupId *int64) (result []structs.ChatChannel, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := getChatChannelsQuery
	if GroupId != nil {
		query = query + " where cc.group_id=" + strconv.FormatInt(*GroupId, 10)
	}
	rows, errQ := pool.Query(query)
	if errQ == nil {
		for rows.Next() {
			var channel structs.ChatChannel
			rows.Scan(&channel.Id, &channel.Name, &channel.AmperId, &channel.AmperName, &channel.GroupId, &channel.GroupName)
			result = append(result, channel)
		}
	} else {
		err = errors.New("unable to run query against database to get chat channels")
		log.Print(errQ.Error(), errQ)
	}
	rows.Close()
	return result, err
}

const getUserChatChannelsQuery = "SELECT cc.id as id, cc.name as name, cc.amper_id as amperId, a.name as amperName, cc.group_id as groupId, ccg.name as groupName, cc.batch_ids as batchIds FROM amper.chat_channel_sys as cc inner join amper.amper_sys as a on a.id = cc.amper_id inner join amper.chat_channel_group_sys as ccg on ccg.id = cc.group_id"

func FetchUserChatChannels(userID *int64) (result []structs.ChatChannel, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := getUserChatChannelsQuery + " WHERE cc.user_ids LIKE '%\\_" + strconv.FormatInt(*userID, 10) + "\\_%'"

	rows, errQ := pool.Query(query)
	if errQ == nil {
		for rows.Next() {
			var channel structs.ChatChannel
			rows.Scan(&channel.Id, &channel.Name, &channel.AmperId, &channel.AmperName, &channel.GroupId, &channel.GroupName, &channel.BatchIds)
			result = append(result, channel)
		}
	} else {
		err = errors.New("unable to run query against database to get chat channels")
		log.Print(errQ.Error(), errQ)
	}
	rows.Close()
	return result, err
}

const getChatChannelQuery = "SELECT * FROM amper.chat_channel_sys where id = %d"

func FetchChatChannel(ChannelId *int64) (result *structs.ChatChannel, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(getChatChannelQuery, *ChannelId)
	result = &structs.ChatChannel{}
	errQR := pool.QueryRow(query).Scan(&result.Id, &result.Name, &result.GroupId, &result.AmperId, &result.UserIds, &result.BatchIds)
	if errQR != nil {
		err = errors.New("unable to run query against database to get chat channel")
		util.Loggify(errQR)
		return nil, err
	}
	return result, nil
}

const updateChatChannelQuery = "UPDATE amper.chat_channel_sys SET user_ids='%s' where id = %d"

func UpdateChatChannel(userID *int64, channel *structs.ChatChannel) (result bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(updateChatChannelQuery, *channel.UserIds, *channel.Id)
	_, errU := pool.Exec(query)
	if errU != nil {
		err = errors.New("unable to update the channel and add new users")
		util.Loggify(errU)
		return false, err
	}
	return true, nil
}

const updateChatChannelBatchIdsQuery = "UPDATE amper.chat_channel_sys SET batch_ids='%s' where id = %d"

func UpdateChatChannelBatchIds(userID *int64, ChannelId int64, BatchIds *string) (result bool, err error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(updateChatChannelBatchIdsQuery, *BatchIds, ChannelId)
	_, errU := pool.Exec(query)
	if errU != nil {
		err = errors.New("unable to update the channel batch ids")
		util.Loggify(errU)
		return false, err
	}
	return true, nil
}

const removeChatChannelGroup = "DELETE from amper.chat_channel_group_sys where id=%d"

func RemoveChatChannelGroup(userID *int64, groupId *int64) (bool, error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(removeChatChannelGroup, *groupId)
	_, errDB := pool.Exec(query)
	if errDB != nil {
		log.Println(errDB.Error(), errDB)
		return false, fmt.Errorf("unable to remove chat channel group with id %d", *groupId)
	}
	return true, nil
}

const removeChatChannel = "DELETE from amper.chat_channel_sys where id=%d"

func RemoveChatChannel(userID *int64, channelId *int64) (bool, error) {
	var pool *sql.DB = databasecache.Pool()
	query := fmt.Sprintf(removeChatChannel, *channelId)
	_, errDB := pool.Exec(query)
	if errDB != nil {
		log.Println(errDB.Error(), errDB)
		return false, fmt.Errorf("unable to remove chat channel with id %d", *channelId)
	}
	return true, nil
}
