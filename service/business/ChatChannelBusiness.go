package business

import (
	"amper/cache/business"
	"amper/common/structs"
	"amper/common/util"
	"amper/common/util/arrays"
	"amper/common/util/datetime"
	"amper/data/database"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	csmap "github.com/gnazaryan/concurrent-swiss-map"
)

var CHAT_CHANNEL_PATH = filepath.Join("__system__", "Chat", "Channel")

var mutexTimoutDuration time.Duration = 30 * time.Minute

func ClearChannelMutexMaps() {
	timedoutKeys := []string{}
	batchUpdateTimedMutexMap.Range(func(key string, value *structs.TimedMutex) (stop bool) {
		if value != nil {
			if time.Now().After(value.Time.Add(mutexTimoutDuration)) {
				timedoutKeys = append(timedoutKeys, key)
			}
		}
		return false
	})
	for _, key := range timedoutKeys {
		batchUpdateTimedMutexMap.Delete(key)
	}
}

var channelReplyMutexTimoutDuration time.Duration = 30 * time.Minute

func ClearChannelReplyMutexMaps() {
	timedoutKeys := []string{}
	batchReplyUpdateTimedMutexMap.Range(func(key string, value *structs.TimedMutex) (stop bool) {
		if value != nil {
			if time.Now().After(value.Time.Add(channelReplyMutexTimoutDuration)) {
				timedoutKeys = append(timedoutKeys, key)
			}
		}
		return false
	})
	for _, key := range timedoutKeys {
		batchReplyUpdateTimedMutexMap.Delete(key)
	}
}

var cacheTimoutDuration time.Duration = 30 * time.Minute

func ClearChannelTimedoutMessages() {
	timedoutKeys := []string{}
	batchMessagesTimedCache.Range(func(key string, value *structs.TimedMessages) (stop bool) {
		if value != nil {
			if time.Now().After(value.Time.Add(cacheTimoutDuration)) {
				timedoutKeys = append(timedoutKeys, key)
			}
		}
		return false
	})
	for _, key := range timedoutKeys {
		batchMessagesTimedCache.Delete(key)
	}
}

var replyInitBatchTimoutDuration time.Duration = 30 * time.Minute

func ClearChannelReplyMessageInitMap() {
	timedoutKeys := []string{}
	replyMessageTimedMutexMap.Range(func(key string, value *structs.TimedBatchId) (stop bool) {
		if value != nil {
			if time.Now().After(value.Time.Add(replyInitBatchTimoutDuration)) {
				timedoutKeys = append(timedoutKeys, key)
			}
		}
		return false
	})
	for _, key := range timedoutKeys {
		replyMessageTimedMutexMap.Delete(key)
	}
}

var cacheReplyMessagesTimoutDuration time.Duration = 30 * time.Minute

func ClearChannelTimedoutReplyMessages() {
	timedoutKeys := []string{}
	batchReplyMessagesTimedCache.Range(func(key string, value *structs.TimedMessages) (stop bool) {
		if value != nil {
			if time.Now().After(value.Time.Add(cacheReplyMessagesTimoutDuration)) {
				timedoutKeys = append(timedoutKeys, key)
			}
		}
		return false
	})
	for _, key := range timedoutKeys {
		batchReplyMessagesTimedCache.Delete(key)
	}
}

func SaveChannelHistory() {
	newChannelMessages.RangeDelete(func(key string, value []structs.Message) (stop bool) {
		//this mutex map ensures there is no concurrent access for a single batch read and write opperation
		//each mutex belongs to specific batch id
		batchMutex := batchUpdateTimedMutexMap.StoreCompute(key, func(value *structs.TimedMutex) *structs.TimedMutex {
			if value == nil {
				value = &structs.TimedMutex{
					Mutex: &sync.RWMutex{},
					Time:  time.Now(),
				}
			} else {
				value.Time = time.Now()
			}
			return value
		})
		batchMutex.Mutex.Lock()
		defer batchMutex.Mutex.Unlock()
		//delete the timed cache since it is going to be modified
		batchMessagesTimedCache.Delete(key)
		if len(value) > 0 {
			batchIdStruct, errBI := structs.ParseId(&key)
			if errBI != nil {
				util.Loggify(errBI)
				return
			}

			reader, errR := GetFileBody(&batchIdStruct.UserId, &key)
			if errR != nil {
				util.Loggify(errR)
				return false
			}
			existingMessagesData, errPD := io.ReadAll(*reader)
			if errPD != nil {
				util.Loggify(errPD)
				return false
			}
			var existingMessages []structs.Message
			errUM := json.Unmarshal(existingMessagesData, &existingMessages)
			if errUM != nil {
				util.Loggify(errUM)
				return false
			}
			existingMessages = append(existingMessages, value...)
			//since the user may have been set by the message fetch api,
			//it is necessary to set it to nil, to not save additional data
			for i := 0; i < len(existingMessages); i++ {
				existingMessages[i].FromUser = nil
			}
			updatedMessages, errUM := json.Marshal(existingMessages)
			if errUM != nil {
				util.Loggify(errUM)
				return false
			}
			success, errU := Update(&batchIdStruct.UserId, &key, &updatedMessages, util.PointerString(""))
			if errU != nil || !success {
				util.Loggify(errU)
				return false
			}
			//The batch file size exceeded it's maximum size
			//upload a new batch file and append the batch file id
			//to the chat channel batch ids
			if len(updatedMessages) > CHAT_FILE_CHUNK_SIZE {
				var channelId *int64
				if batchIdStruct.OptionalValue != nil {
					optionalValueSplit := strings.Split(*batchIdStruct.OptionalValue, "-")
					if len(optionalValueSplit) == 2 {
						id, errID := strconv.ParseInt(optionalValueSplit[1], 10, 64)
						if errID == nil {
							channelId = &id
						}
					}
				}
				if channelId == nil {
					util.Loggify(fmt.Errorf("the batch id %s doesn't point to a valid channel id", key))
					return false
				}
				channelMutex := channelMutexMap.StoreCompute(*channelId, func(value *sync.RWMutex) *sync.RWMutex {
					if value == nil {
						value = &sync.RWMutex{}
					}
					return value
				})
				channelMutex.Lock()
				defer channelMutex.Unlock()

				messagesBytes := []byte("[]")
				size := int64(len(messagesBytes))
				date := datetime.FormatDate(time.Now())
				optionalValue := util.PointerString("channel-" + strconv.FormatInt(*channelId, 10))
				success, fileMetadata, errU := Upload(&batchIdStruct.UserId, nil, &messagesBytes, util.PointerString(date), util.PointerString("application/json"), &size, util.PointerString(filepath.Join(CHAT_CHANNEL_PATH, strconv.FormatInt(*channelId, 10))), optionalValue)
				if success && errU == nil {
					channel := business.GetChatChannel(channelId)
					if channel != nil {
						chatHistoryItems := []structs.ChatHistoryItem{}
						if channel.BatchIds != nil {
							json.Unmarshal([]byte(*channel.BatchIds), &chatHistoryItems)
							if len(chatHistoryItems) > 0 {
								chatHistoryItems[len(chatHistoryItems)-1].Full = true
								chatHistoryItems = append(chatHistoryItems, structs.ChatHistoryItem{
									Id:   fileMetadata.Id,
									Full: false,
								})
								chatHistoryItemsBytes, errCHIB := json.Marshal(chatHistoryItems)
								if errCHIB == nil {
									BatchIds := string(chatHistoryItemsBytes)
									instance := business.GetAmperInstance(*business.AmperId())
									sessionId := GenerateSessionId(instance.Identifier, instance.Key, util.PointerString("app"))
									_, errUCCB := UpdateChatChannelBatchIds(&batchIdStruct.UserId, &sessionId, *channelId, &BatchIds)
									if errUCCB != nil {
										util.Loggify(errUCCB)
									} else {
										//Send update to all active participants for the new batch added
										//All active ui users should add the batch id to the chat item in client side
										userMessageUpdate := structs.UserMessageUpdate{
											MessageType:    util.PointerString("channel"),
											UpdateType:     util.PointerString("newBatch"),
											OpperationType: util.PointerString("newBatch"),
											From:           util.PointerString(strconv.FormatInt(SYSTEM_USER_ID, 10)),
											To:             util.PointerString(strconv.FormatInt(*channelId, 10)),
											Value:          fileMetadata.Id,
										}
										//send update to all participants with a defer to not bloeck the current thread
										defer sendUpdateToChannelParticipants(&SYSTEM_USER_ID, &sessionId, channel, &userMessageUpdate, "chat")
									}
								}
							} else {
								util.Loggify(fmt.Errorf("the batch ids of the channel %d doesn't contain any valid batch id, the channel is not configured properly", *channelId))
							}
						} else {
							util.Loggify(fmt.Errorf("the batch ids of the channel %d is empty, the channel is not configured properly", *channelId))
						}
					} else {
						util.Loggify(fmt.Errorf("not able to retrieve the chat channel with id %d from cache", *channelId))
						return false
					}
				}
			}
		}
		return false
	})
}

func SaveChannelMessageUpdates() {
	channelMessageUpdatesCS.RangeDelete(func(key string, value map[string][]structs.MessageUpdate) (stop bool) {
		//this mutex map ensures there is no concurrent access for a single batch read and write opperation
		//each mutex belongs to specific batch id
		batchMutex := batchUpdateTimedMutexMap.StoreCompute(key, func(value *structs.TimedMutex) *structs.TimedMutex {
			if value == nil {
				value = &structs.TimedMutex{
					Mutex: &sync.RWMutex{},
					Time:  time.Now(),
				}
			} else {
				value.Time = time.Now()
			}
			return value
		})
		batchMutex.Mutex.Lock()
		defer batchMutex.Mutex.Unlock()
		//delete the timed cache since it is going to be modified
		batchMessagesTimedCache.Delete(key)

		batchIdStruct, errBI := structs.ParseId(&key)
		if errBI != nil {
			util.Loggify(errBI)
			return
		}

		reader, errR := GetFileBody(&batchIdStruct.UserId, &key)
		if errR != nil {
			util.Loggify(errR)
			return false
		}
		existingMessagesData, errPD := io.ReadAll(*reader)
		if errPD != nil {
			util.Loggify(errPD)
			return false
		}
		var existingMessages []structs.Message
		errUM := json.Unmarshal(existingMessagesData, &existingMessages)
		if errUM != nil {
			util.Loggify(errUM)
			return false
		}
		for i := 0; i < len(existingMessages); i++ {
			messageUpdates := value[*existingMessages[i].Id]
			for _, messageUpdate := range messageUpdates {
				fromInt64, errFI := strconv.ParseInt(*messageUpdate.From, 10, 64)
				if errFI != nil {
					util.Loggify(errFI)
					continue
				}
				if *messageUpdate.UpdateType == "reaction" {
					if existingMessages[i].Reactions == nil {
						existingMessages[i].Reactions = map[string][]int64{}
					}
					if *messageUpdate.OpperationType == "add" {
						if existingMessages[i].Reactions[*messageUpdate.Value] == nil {
							existingMessages[i].Reactions[*messageUpdate.Value] = make([]int64, 0)
							existingMessages[i].Reactions[*messageUpdate.Value] = append(existingMessages[i].Reactions[*messageUpdate.Value], fromInt64)
						} else {
							if !arrays.Contains(existingMessages[i].Reactions[*messageUpdate.Value], fromInt64) {
								existingMessages[i].Reactions[*messageUpdate.Value] = append(existingMessages[i].Reactions[*messageUpdate.Value], fromInt64)
							}
						}
					} else if *messageUpdate.OpperationType == "remove" {
						if existingMessages[i].Reactions[*messageUpdate.Value] != nil {
							existingMessages[i].Reactions[*messageUpdate.Value] = arrays.Remove(existingMessages[i].Reactions[*messageUpdate.Value], fromInt64)
						}
					}
				} else if *messageUpdate.UpdateType == "edit" {
					existingMessages[i].Text = messageUpdate.Value
				} else if *messageUpdate.UpdateType == "remove" {
					existingMessages[i].Deleted = true
				} else if *messageUpdate.UpdateType == "reply" {
					if existingMessages[i].Replies == nil {
						existingMessages[i].Replies = make(map[string]int)
					}
					existingMessages[i].Replies[*messageUpdate.Value]++
				} else if *messageUpdate.UpdateType == "replyInitialisation" {
					existingMessages[i].ReplyBatchId = messageUpdate.Value
				}
			}
		}
		updatedMessages, errUM := json.Marshal(existingMessages)
		if errUM != nil {
			util.Loggify(errUM)
			return false
		}
		success, errU := Update(&batchIdStruct.UserId, &key, &updatedMessages, util.PointerString(""))
		if errU != nil || !success {
			util.Loggify(errU)
			return false
		}
		return false
	})
}

func SaveChannelRepliesHistory() {
	newChannelRepliesCS.RangeDelete(func(key string, value []structs.Message) (stop bool) {
		//this mutex map ensures there is no concurrent access for a single batch read and write opperation
		//each mutex belongs to specific batch id
		batchReplyMutex := batchReplyUpdateTimedMutexMap.StoreCompute(key, func(value *structs.TimedMutex) *structs.TimedMutex {
			if value == nil {
				value = &structs.TimedMutex{
					Mutex: &sync.RWMutex{},
					Time:  time.Now(),
				}
			} else {
				value.Time = time.Now()
			}
			return value
		})
		batchReplyMutex.Mutex.Lock()
		defer batchReplyMutex.Mutex.Unlock()
		//delete the messages cache, since it is going to be updated soon
		batchReplyMessagesTimedCache.Delete(key)

		batchIdStruct, errBI := structs.ParseId(&key)
		if errBI != nil {
			util.Loggify(errBI)
			return
		}

		reader, errR := GetFileBody(&batchIdStruct.UserId, &key)
		if errR != nil {
			util.Loggify(errR)
			return false
		}
		existingMessagesData, errPD := io.ReadAll(*reader)
		if errPD != nil {
			util.Loggify(errPD)
			return false
		}
		var existingMessages []structs.Message
		errUM := json.Unmarshal(existingMessagesData, &existingMessages)
		if errUM != nil {
			util.Loggify(errUM)
			return false
		}
		existingMessages = append(existingMessages, value...)
		for i := 0; i < len(existingMessages); i++ {
			existingMessages[i].FromUser = nil
		}
		updatedMessages, errUM := json.Marshal(existingMessages)
		if errUM != nil {
			util.Loggify(errUM)
			return false
		}
		success, errU := Update(&batchIdStruct.UserId, &key, &updatedMessages, util.PointerString(""))
		if errU != nil || !success {
			util.Loggify(errU)
			return false
		}
		return false
	})
}

func SaveChannelRepliesMessageUpdates() {
	channelRepliesMessageUpdatesCS.RangeDelete(func(key string, value map[string][]structs.MessageUpdate) (stop bool) {
		//this mutex map ensures there is no concurrent access for a single batch read and write opperation
		//each mutex belongs to specific batch id
		batchReplyMutex := batchReplyUpdateTimedMutexMap.StoreCompute(key, func(value *structs.TimedMutex) *structs.TimedMutex {
			if value == nil {
				value = &structs.TimedMutex{
					Mutex: &sync.RWMutex{},
					Time:  time.Now(),
				}
			} else {
				value.Time = time.Now()
			}
			return value
		})
		batchReplyMutex.Mutex.Lock()
		defer batchReplyMutex.Mutex.Unlock()

		batchReplyMessagesTimedCache.Delete(key)
		batchIdStruct, errBI := structs.ParseId(&key)
		if errBI != nil {
			util.Loggify(errBI)
			return
		}

		reader, errR := GetFileBody(&batchIdStruct.UserId, &key)
		if errR != nil {
			util.Loggify(errR)
			return false
		}
		existingMessagesData, errPD := io.ReadAll(*reader)
		if errPD != nil {
			util.Loggify(errPD)
			return false
		}
		var existingMessages []structs.Message
		errUM := json.Unmarshal(existingMessagesData, &existingMessages)
		if errUM != nil {
			util.Loggify(errUM)
			return false
		}
		for i := 0; i < len(existingMessages); i++ {
			messageUpdates := value[*existingMessages[i].Id]
			for _, messageUpdate := range messageUpdates {
				fromInt64, errFI := strconv.ParseInt(*messageUpdate.From, 10, 64)
				if errFI != nil {
					util.Loggify(errFI)
					continue
				}
				if *messageUpdate.UpdateType == "reaction" {
					if existingMessages[i].Reactions == nil {
						existingMessages[i].Reactions = map[string][]int64{}
					}
					if *messageUpdate.OpperationType == "add" {
						if existingMessages[i].Reactions[*messageUpdate.Value] == nil {
							existingMessages[i].Reactions[*messageUpdate.Value] = make([]int64, 0)
							existingMessages[i].Reactions[*messageUpdate.Value] = append(existingMessages[i].Reactions[*messageUpdate.Value], fromInt64)
						} else {
							if !arrays.Contains(existingMessages[i].Reactions[*messageUpdate.Value], fromInt64) {
								existingMessages[i].Reactions[*messageUpdate.Value] = append(existingMessages[i].Reactions[*messageUpdate.Value], fromInt64)
							}
						}
					} else if *messageUpdate.OpperationType == "remove" {
						if existingMessages[i].Reactions[*messageUpdate.Value] != nil {
							existingMessages[i].Reactions[*messageUpdate.Value] = arrays.Remove(existingMessages[i].Reactions[*messageUpdate.Value], fromInt64)
						}
					}
				} else if *messageUpdate.UpdateType == "edit" {
					existingMessages[i].Text = messageUpdate.Value
				} else if *messageUpdate.UpdateType == "remove" {
					existingMessages[i].Deleted = true
				} else if *messageUpdate.UpdateType == "reply" {
					if existingMessages[i].Replies == nil {
						existingMessages[i].Replies = make(map[string]int)
					}
					existingMessages[i].Replies[*messageUpdate.Value]++
				} else if *messageUpdate.UpdateType == "replyInitialisation" {
					existingMessages[i].ReplyBatchId = messageUpdate.Value
				}
			}
		}
		updatedMessages, errUM := json.Marshal(existingMessages)
		if errUM != nil {
			util.Loggify(errUM)
			return false
		}
		success, errU := Update(&batchIdStruct.UserId, &key, &updatedMessages, util.PointerString(""))
		if errU != nil || !success {
			util.Loggify(errU)
			return false
		}
		return false
	})
}

func GetChatChannelState(userID *int64, sessionId *string) (result *[]structs.ChatChannelsGroup, err error) {
	intermediateResult := []structs.ChatChannelsGroup{}
	userChannels, errUC := database.FetchUserChatChannels(userID)
	var channelGroupMap map[int64][]structs.ChatChannel = make(map[int64][]structs.ChatChannel)
	if errUC == nil {
		for _, userChannel := range userChannels {
			if channelGroupMap[*userChannel.GroupId] == nil {
				channelGroupMap[*userChannel.GroupId] = []structs.ChatChannel{
					userChannel,
				}
			} else {
				channelGroupMap[*userChannel.GroupId] = append(channelGroupMap[*userChannel.GroupId], userChannel)
			}
		}
	}
	for groupId, userChannels := range channelGroupMap {
		GroupName := ""
		Channels := []structs.ChatChannelItem{}
		for _, userChannel := range userChannels {
			GroupName = *userChannel.GroupName
			if userChannel.BatchIds == nil {
				continue
			}
			chatHistoryItems := []structs.ChatHistoryItem{}
			json.Unmarshal([]byte(*userChannel.BatchIds), &chatHistoryItems)
			to := strconv.FormatInt(*userChannel.Id, 10)
			Channels = append(Channels, structs.ChatChannelItem{
				ChannelId: *userChannel.Id,
				Label:     userChannel.Name,
				AmperId:   *userChannel.AmperId,
				ChatHistory: &structs.ChatHistory{
					To:           &to,
					HistoryItems: chatHistoryItems,
				},
			})
		}
		chatChannelGroup := structs.ChatChannelsGroup{
			GroupId:  groupId,
			Label:    &GroupName,
			Channels: &Channels,
		}
		intermediateResult = append(intermediateResult, chatChannelGroup)
	}
	sort.Slice(intermediateResult, func(i, j int) bool {
		if intermediateResult[i].Label == nil {
			return false
		}
		if intermediateResult[j].Label == nil {
			return true
		}
		return (strings.Compare(*intermediateResult[i].Label, *intermediateResult[j].Label) < 0)
	})
	return &intermediateResult, nil
}

func CreateChatGroup(userID *int64, Name *string) (success bool, err error) {
	return database.CreateChatGroup(userID, Name)
}

func FetchChatChannelGroups(userID *int64) (result []structs.ChatChannelGroup, err error) {
	return database.FetchChatChannelGroups(userID)
}

func CreateChatChannel(userID *int64, Name *string, AmperId *int64, GroupId *int64) (success bool, err error) {
	channelId, errCCC := database.CreateChatChannel(userID, Name, AmperId, GroupId)
	if errCCC == nil {
		messagesBytes := []byte("[]")
		size := int64(len(messagesBytes))
		date := datetime.FormatDate(time.Now())
		optionalValue := util.PointerString("channel-" + strconv.FormatInt(channelId, 10))
		successU, fileMetadata, errU := Upload(userID, nil, &messagesBytes, util.PointerString(date), util.PointerString("application/json"), &size, util.PointerString(filepath.Join(CHAT_CHANNEL_PATH, strconv.FormatInt(channelId, 10))), optionalValue)
		if successU && errU == nil {
			chatHistoryItems := []structs.ChatHistoryItem{{
				Id:   fileMetadata.Id,
				Full: false,
			}}
			chatHistoryItemsBytes, errCHIB := json.Marshal(chatHistoryItems)
			if errCHIB == nil {
				BatchIds := string(chatHistoryItemsBytes)
				return database.UpdateChatChannelBatchIds(userID, channelId, &BatchIds)
			}
		}
	} else {
		util.Loggify(errCCC)
	}

	return false, fmt.Errorf("not able to craete a channel this time, try again later or contact the support")
}

func UpdateChatChannelBatchIds(userID *int64, sessionId *string, ChannelId int64, BatchIds *string) (result bool, err error) {
	successUCCB, errUCCB := database.UpdateChatChannelBatchIds(userID, ChannelId, BatchIds)
	if successUCCB && errUCCB == nil {
		result = true
		instanceTypes := []string{"amperInstance"}
		successFC, errFC := FederatedCall(userID, sessionId, map[string]string{
			"amperDatastoreInstance": "amper-datastore/invalidateCache",
			"amperInstance":          "amper/invalidateCache",
		}, map[string]interface{}{
			"ChatChannelId": strconv.FormatInt(ChannelId, 10),
			"name":          "chatChannel",
		}, &instanceTypes)
		util.Loggify(errFC)
		if !successFC || errFC != nil {
			err = fmt.Errorf("not able to reset the federated cache for chat channel, contact the support")
			result = false
		}
		//Since channels are dedicated to specific amper instance, it is
		//expected that the callee of this function will be on the amper instance
		//where this cache will be valid
		business.InvalidateChatChannel(&ChannelId)
		result = true
	} else {
		util.Loggify(errUCCB)
		err = fmt.Errorf("not able to update the batch ids for channel %d", ChannelId)
	}
	return result, err
}

func FetchChatChannels(userID *int64, GroupId *int64) (result []structs.ChatChannel, err error) {
	return database.FetchChatChannels(userID, GroupId)
}

func RemoveChatChannelGroup(userID *int64, groupId *int64) (bool, error) {
	return database.RemoveChatChannelGroup(userID, groupId)
}

func RemoveChatChannel(userID *int64, channelId *int64) (bool, error) {
	return database.RemoveChatChannel(userID, channelId)
}

var addUserToChannelMutex sync.Mutex

func AddUsersToChannel(userID *int64, sessionId *string, ChannelId *int64, UserIds *[]int64) (success bool, err error) {
	addUserToChannelMutex.Lock()
	defer addUserToChannelMutex.Unlock()
	channel, errC := database.FetchChatChannel(ChannelId)
	if errC != nil {
		util.Loggify(errC)
		return false, fmt.Errorf("not able to locate a channel with the given channel Id")
	}
	if channel.UserIds == nil {
		channel.UserIds = util.PointerString("")
	}
	for _, userId := range *UserIds {
		userIdWrapped := "_" + strconv.FormatInt(userId, 10) + "_"
		if !strings.Contains(*channel.UserIds, userIdWrapped) {
			channel.UserIds = util.PointerString(*channel.UserIds + userIdWrapped)
		}
	}
	successUpdate, errSU := database.UpdateChatChannel(userID, channel)
	if successUpdate && errSU == nil {
		success = true
		instanceTypes := []string{"amperInstance"}
		successFC, errFC := FederatedCall(userID, sessionId, map[string]string{
			"amperDatastoreInstance": "amper-datastore/invalidateCache",
			"amperInstance":          "amper/invalidateCache",
		}, map[string]interface{}{
			"ChatChannelId": strconv.FormatInt(*ChannelId, 10),
			"name":          "chatChannel",
		}, &instanceTypes)
		util.Loggify(errFC)
		if !successFC || errFC != nil {
			err = fmt.Errorf("not able to reset the federated cache for chat channel, contact the support")
			success = false
		}
	}
	util.Loggify(errSU)
	return success, err
}

func FetchChatChannelUsers(userID *int64, ChannelId *int64, Search *[]string, Start int, Limit int) (result *[]structs.User, resultTotalCount int, err error) {
	channel, errC := database.FetchChatChannel(ChannelId)
	if errC == nil {
		if len(*channel.UserIds) > 0 {
			channelUsersSplit := strings.Split(*channel.UserIds, "__")
			var userIds []string = make([]string, 0)
			for _, channelUserDirty := range channelUsersSplit {
				userIds = append(userIds, strings.Replace(channelUserDirty, "_", "", -1))
			}
			users, totalCount, errU := database.GetUsersIn(&Start, &Limit, Search, userIds)
			if errU != nil {
				util.Loggify(errU)
				err = fmt.Errorf("not able to fetch channel users, try again later or contact the support")
			} else {
				result = &users
				resultTotalCount = totalCount
			}
		} else {
			result = &[]structs.User{}
		}
	} else {
		util.Loggify(errC)
	}
	return result, resultTotalCount, err
}

func RemoveChatChannelUser(userID *int64, sessionId *string, ChannelId *int64, UserIds *[]int64) (success bool, err error) {
	addUserToChannelMutex.Lock()
	defer addUserToChannelMutex.Unlock()
	channel, errC := database.FetchChatChannel(ChannelId)
	if errC != nil {
		util.Loggify(errC)
		return false, fmt.Errorf("not able to locate a channel with the given channel Id")
	}
	if channel.UserIds == nil {
		channel.UserIds = util.PointerString("")
	}
	for _, userId := range *UserIds {
		userIdWrapped := "_" + strconv.FormatInt(userId, 10) + "_"
		if strings.Contains(*channel.UserIds, userIdWrapped) {
			channel.UserIds = util.PointerString(strings.Replace(*channel.UserIds, userIdWrapped, "", -1))
		}
	}
	successUpdate, errSU := database.UpdateChatChannel(userID, channel)
	if successUpdate && errSU == nil {
		success = true
		instanceTypes := []string{"amperInstance"}
		successFC, errFC := FederatedCall(userID, sessionId, map[string]string{
			"amperDatastoreInstance": "amper-datastore/invalidateCache",
			"amperInstance":          "amper/invalidateCache",
		}, map[string]interface{}{
			"ChatChannelId": strconv.FormatInt(*ChannelId, 10),
			"name":          "chatChannel",
		}, &instanceTypes)
		util.Loggify(errFC)
		if !successFC || errFC != nil {
			err = fmt.Errorf("not able to reset the federated cache for chat channel, contact the support")
			success = false
		}
	}
	util.Loggify(errSU)
	return success, err
}

var channelMutexMap = csmap.Create[int64, *sync.RWMutex](
	// set the number of map shards. the default value is 32.
	csmap.WithShardCount[int64, *sync.RWMutex](100),

	// if don't set custom hasher, use the built-in maphash.
	csmap.WithCustomHasher[int64, *sync.RWMutex](func(key int64) uint64 {
		return uint64(key)
	}),

	// set the total capacity, every shard map has total capacity/shard count capacity. the default value is 0.
	csmap.WithSize[int64, *sync.RWMutex](10000))

var newChannelMessages = csmap.Create[string, []structs.Message](
	// set the number of map shards. the default value is 32.
	csmap.WithShardCount[string, []structs.Message](100),

	// if don't set custom hasher, use the built-in maphash.
	csmap.WithCustomHasher[string, []structs.Message](func(key string) uint64 {
		hash := fnv.New64a()
		hash.Write([]byte(key))
		return hash.Sum64()
	}),

	// set the total capacity, every shard map has total capacity/shard count capacity. the default value is 0.
	csmap.WithSize[string, []structs.Message](10000))

func sendChannelMessage(userID *int64, sessionId *string, from *string, to *string, Type *string, id *string, text *string, Attachments *[]structs.ChatAttachment) (success bool, message *structs.Message, err error) {
	toInt64, errTI := strconv.ParseInt(*to, 10, 64)
	if errTI != nil {
		return false, nil, fmt.Errorf("sendChannelMessage - the 'to' parameter supplied doesn't resolve to a valid identifying number")
	}
	fromInt64, errFI := strconv.ParseInt(*from, 10, 64)
	if errFI != nil {
		return false, nil, fmt.Errorf("sendChannelMessage - the 'from' parameter supplied doesn't resolve to a valid identifying number")
	}
	channelMutex := channelMutexMap.StoreCompute(toInt64, func(value *sync.RWMutex) *sync.RWMutex {
		if value == nil {
			value = &sync.RWMutex{}
		}
		return value
	})

	//acquire a channel lock, since the batch id can be added synchroniously by the chat history saving thread
	channelMutex.RLock()
	channel := business.GetChatChannel(&toInt64)
	if channel == nil {
		channelMutex.RUnlock()
		return false, nil, fmt.Errorf("the to parameter supplied doesn't resolve to a valid existing channel")
	}
	if len(channel.BatchIdsArray) < 1 {
		channelMutex.RUnlock()
		return false, nil, fmt.Errorf("the specified channel `%s` is not properly configured", *to)
	}
	if !strings.Contains(*channel.UserIds, "_"+strconv.FormatInt(*userID, 10)+"_") {
		channelMutex.RUnlock()
		return false, nil, fmt.Errorf("the user doesn't belong to the specified channel")
	}
	BatchId := channel.BatchIdsArray[len(channel.BatchIdsArray)-1].Id
	channelMutex.RUnlock()

	messageTime := time.Now().UnixMilli()
	message = &structs.Message{
		From:        from,
		FromUser:    business.GetUser(&fromInt64, true),
		To:          to,
		Id:          id,
		Text:        text,
		DateTime:    messageTime,
		BatchId:     BatchId,
		Attachments: Attachments,
	}

	newChannelMessages.StoreCompute(*BatchId, func(value []structs.Message) []structs.Message {
		if value != nil {
			return append(value, *message)
		} else {
			return []structs.Message{*message}
		}
	})
	userMessageUpdate := structs.UserMessageUpdate{
		Message:        message,
		MessageType:    util.PointerString("channel"),
		UpdateType:     util.PointerString("newMessage"),
		OpperationType: util.PointerString("newMessage"),
		From:           from,
		To:             to,
	}
	//send update to all participants with a defer to not bloeck the current thread
	defer sendUpdateToChannelParticipants(userID, sessionId, channel, &userMessageUpdate, "chat")
	return true, message, nil
}

func sendUpdateToChannelParticipants(userID *int64, sessionId *string, channel *structs.ChatChannel, messageUpdate *structs.UserMessageUpdate, category string) {
	for instanceId, userIds := range channel.InstanceUserIDs {
		selfExludedParticipants := []int64{}
		for _, participant := range userIds {
			if *userID != participant {
				selfExludedParticipants = append(selfExludedParticipants, participant)
			}
		}
		instance := business.GetAmperInstance(instanceId)
		_, _, errPS := DedicatedCallWithRetry(userID, sessionId, map[string]string{
			"amperInstance": "updates/push",
		}, map[string]interface{}{
			"category":     category,
			"participants": selfExludedParticipants,
			"value":        messageUpdate,
		}, instance)
		if errPS != nil {
			util.Loggify(errPS)
		}
	}
}

// this mutex map ensures there is no concurrent access for a single batch read and write opperation
// each mutex belongs to specific batch id
var batchUpdateTimedMutexMap = csmap.Create[string, *structs.TimedMutex](
	// set the number of map shards. the default value is 32.
	csmap.WithShardCount[string, *structs.TimedMutex](100),

	// if don't set custom hasher, use the built-in maphash.
	csmap.WithCustomHasher[string, *structs.TimedMutex](func(key string) uint64 {
		hash := fnv.New64a()
		hash.Write([]byte(key))
		return hash.Sum64()
	}),

	// set the total capacity, every shard map has total capacity/shard count capacity. the default value is 0.
	csmap.WithSize[string, *structs.TimedMutex](10000))

var batchMessagesTimedCache = csmap.Create[string, *structs.TimedMessages](
	// set the number of map shards. the default value is 32.
	csmap.WithShardCount[string, *structs.TimedMessages](100),

	// if don't set custom hasher, use the built-in maphash.
	csmap.WithCustomHasher[string, *structs.TimedMessages](func(key string) uint64 {
		hash := fnv.New64a()
		hash.Write([]byte(key))
		return hash.Sum64()
	}),

	// set the total capacity, every shard map has total capacity/shard count capacity. the default value is 0.
	csmap.WithSize[string, *structs.TimedMessages](10000))

func fetchChatChannel(userID *int64, from *string, to *string, id *string, includeLatest bool) (result []structs.Message, participants map[int64]*structs.User, err error) {
	// this mutex map ensures there is no concurrent access for a single batch read and write opperation
	// each mutex belongs to specific batch id
	batchMutex := batchUpdateTimedMutexMap.StoreCompute(*id, func(value *structs.TimedMutex) *structs.TimedMutex {
		if value == nil {
			value = &structs.TimedMutex{
				Mutex: &sync.RWMutex{},
				Time:  time.Now(),
			}
		} else {
			value.Time = time.Now()
		}
		return value
	})
	batchMutex.Mutex.RLock()
	defer batchMutex.Mutex.RUnlock()

	participants = make(map[int64]*structs.User)
	result = make([]structs.Message, 0)
	timedMessages := batchMessagesTimedCache.StoreCompute(*id, func(value *structs.TimedMessages) *structs.TimedMessages {
		if value == nil {
			reader, errR := GetFileBody(userID, id)
			if errR != nil {
				util.Loggify(errR)
				return nil
			}
			existingMessagesData, errPD := io.ReadAll(*reader)
			if errPD != nil {
				util.Loggify(errPD)
				return nil
			}
			var existingMessages []structs.Message
			errUM := json.Unmarshal(existingMessagesData, &existingMessages)
			if errUM != nil {
				util.Loggify(errUM)
				return nil
			}

			value = &structs.TimedMessages{
				Time:     time.Now(),
				Messages: &existingMessages,
			}
		} else {
			value.Time = time.Now()
		}
		return value
	})
	var intermediateResult []structs.Message = make([]structs.Message, 0)
	for _, existingMessage := range *timedMessages.Messages {
		if !existingMessage.Deleted {
			intermediateResult = append(intermediateResult, existingMessage)
		}
	}
	if includeLatest {
		newMessages, ok := newChannelMessages.Load(*id)
		if ok {
			intermediateResult = append(intermediateResult, newMessages...)
		}
	}

	channelMessageUpdatesCS.LoadLocked(*id, func(value map[string][]structs.MessageUpdate) {
	FirstLoop:
		for i := 0; i < len(intermediateResult); i++ {
			messageUpdates := value[*intermediateResult[i].Id]
			for _, messageUpdate := range messageUpdates {
				fromInt64, errFI := strconv.ParseInt(*messageUpdate.From, 10, 64)
				if errFI != nil {
					util.Loggify(errFI)
					continue
				}
				if *messageUpdate.UpdateType == "reaction" {
					if intermediateResult[i].Reactions == nil {
						intermediateResult[i].Reactions = map[string][]int64{}
					}
					if *messageUpdate.OpperationType == "add" {
						if intermediateResult[i].Reactions[*messageUpdate.Value] == nil {
							intermediateResult[i].Reactions[*messageUpdate.Value] = make([]int64, 0)
							intermediateResult[i].Reactions[*messageUpdate.Value] = append(intermediateResult[i].Reactions[*messageUpdate.Value], fromInt64)
						} else {
							if !arrays.Contains(intermediateResult[i].Reactions[*messageUpdate.Value], fromInt64) {
								intermediateResult[i].Reactions[*messageUpdate.Value] = append(intermediateResult[i].Reactions[*messageUpdate.Value], fromInt64)
							}
						}
					} else if *messageUpdate.OpperationType == "remove" {
						if intermediateResult[i].Reactions[*messageUpdate.Value] != nil {
							intermediateResult[i].Reactions[*messageUpdate.Value] = arrays.Remove(intermediateResult[i].Reactions[*messageUpdate.Value], fromInt64)
						}
					}
				} else if *messageUpdate.UpdateType == "edit" {
					intermediateResult[i].Text = messageUpdate.Value
				} else if *messageUpdate.UpdateType == "remove" {
					continue FirstLoop
				} else if *messageUpdate.UpdateType == "reply" {
					if intermediateResult[i].Replies == nil {
						intermediateResult[i].Replies = make(map[string]int)
					}
					intermediateResult[i].Replies[*messageUpdate.Value]++
				} else if *messageUpdate.UpdateType == "replyInitialisation" {
					intermediateResult[i].ReplyBatchId = messageUpdate.Value
				}
			}
			//add the from user if missing in the participants map
			if intermediateResult[i].From != nil {
				fromUserInt64, _ := strconv.ParseInt(*intermediateResult[i].From, 10, 64)
				if participants[fromUserInt64] == nil {
					participants[fromUserInt64] = business.GetUser(&fromUserInt64, true)
				}
			}
			//add users reacting to the message to the participants map
			if intermediateResult[i].Reactions != nil {
				for _, userIds := range intermediateResult[i].Reactions {
					if len(userIds) > 0 {
						for _, userId := range userIds {
							if participants[userId] == nil {
								participants[userId] = business.GetUser(&userId, true)
							}
						}
					}
				}
			}
			result = append(result, intermediateResult[i])
		}
	})
	return result, participants, nil
}

var channelMessageUpdatesCS = csmap.Create[string, map[string][]structs.MessageUpdate](
	// set the number of map shards. the default value is 32.
	csmap.WithShardCount[string, map[string][]structs.MessageUpdate](100),

	// if don't set custom hasher, use the built-in maphash.
	csmap.WithCustomHasher[string, map[string][]structs.MessageUpdate](func(key string) uint64 {
		hash := fnv.New64a()
		hash.Write([]byte(key))
		return hash.Sum64()
	}),

	// set the total capacity, every shard map has total capacity/shard count capacity. the default value is 0.
	csmap.WithSize[string, map[string][]structs.MessageUpdate](10000))

func updateReceiveChatChannel(userID *int64, sessionId *string, BatchId *string, MessageId *string, MessageType *string, UpdateType *string, OpperationType *string, Value *string, From *string, To *string) (success bool, err error) {
	toInt64, errTI := strconv.ParseInt(*To, 10, 64)
	if errTI != nil {
		return false, fmt.Errorf("updateReceiveChatChannel - the 'to' parameter supplied doesn't resolve to a valid identifying number")
	}
	fromInt64, errFI := strconv.ParseInt(*From, 10, 64)
	if errFI != nil {
		return false, fmt.Errorf("updateReceiveChatChannel - the 'to' parameter supplied doesn't resolve to a valid identifying number")
	}
	messageUpdate := structs.MessageUpdate{
		MessageType:    MessageType,
		UpdateType:     UpdateType,
		OpperationType: OpperationType,
		From:           From,
		To:             To,
		Value:          Value,
	}

	channelMessageUpdatesCS.StoreCompute(*BatchId, func(value map[string][]structs.MessageUpdate) map[string][]structs.MessageUpdate {
		if value != nil {
			if value[*MessageId] != nil {
				value[*MessageId] = append(value[*MessageId], messageUpdate)
			} else {
				value[*MessageId] = []structs.MessageUpdate{messageUpdate}
			}
		} else {
			value = make(map[string][]structs.MessageUpdate)
			value[*MessageId] = []structs.MessageUpdate{messageUpdate}
		}
		return value
	})
	channel := business.GetChatChannel(&toInt64)
	if channel == nil {
		return false, fmt.Errorf("the to parameter supplied doesn't resolve to a valid existing channel")
	}
	userMessageUpdate := structs.UserMessageUpdate{
		Message: &structs.Message{
			Id:       MessageId,
			FromUser: business.GetUser(&fromInt64, true),
		},
		MessageType:    MessageType,
		UpdateType:     UpdateType,
		OpperationType: OpperationType,
		From:           From,
		To:             To,
		Value:          Value,
	}
	//send update to all participants with a defer to not bloeck the current thread
	defer sendUpdateToChannelParticipants(userID, sessionId, channel, &userMessageUpdate, "chat")
	return true, nil
}

// this mutex map ensures there is no concurrent access for a single reply to a message
// since three can possibly be many replies to the same message simultainusley, should be made sure
// there are no simultainuse replies overwriting each other
var replyMessageTimedMutexMap = csmap.Create[string, *structs.TimedBatchId](
	// set the number of map shards. the default value is 32.
	csmap.WithShardCount[string, *structs.TimedBatchId](100),

	// if don't set custom hasher, use the built-in maphash.
	csmap.WithCustomHasher[string, *structs.TimedBatchId](func(key string) uint64 {
		hash := fnv.New64a()
		hash.Write([]byte(key))
		return hash.Sum64()
	}),

	// set the total capacity, every shard map has total capacity/shard count capacity. the default value is 0.
	csmap.WithSize[string, *structs.TimedBatchId](10000))

var newChannelRepliesCS = csmap.Create[string, []structs.Message](
	// set the number of map shards. the default value is 32.
	csmap.WithShardCount[string, []structs.Message](100),

	// if don't set custom hasher, use the built-in maphash.
	csmap.WithCustomHasher[string, []structs.Message](func(key string) uint64 {
		hash := fnv.New64a()
		hash.Write([]byte(key))
		return hash.Sum64()
	}),

	// set the total capacity, every shard map has total capacity/shard count capacity. the default value is 0.
	csmap.WithSize[string, []structs.Message](10000))

func sendChannelReplyMessage(userID *int64, sessionId *string, BatchId *string, From *string, To *string, RepliesToMessageId *string, RepliesToMessageBatchId *string, RepliesToMessageBatchId1 *string, RepliesToMessageType *string, Id *string, Text *string, Attachments *[]structs.ChatAttachment, Participants *[]string) (success bool, message *structs.Message, err error) {
	channelId, errTI := strconv.ParseInt(*To, 10, 64)
	if errTI != nil {
		return false, nil, fmt.Errorf("sendChannelReplyMessage - the 'to' parameter supplied doesn't resolve to a valid identifying number")
	}
	fromInt64, errFI := strconv.ParseInt(*From, 10, 64)
	if errFI != nil {
		return false, nil, fmt.Errorf("sendChannelReplyMessage - the 'from' parameter supplied doesn't resolve to a valid identifying number")
	}
	channel := business.GetChatChannel(&channelId)
	if *channel.AmperId != *business.AmperId() {
		return false, nil, fmt.Errorf("sendChannelReplyMessage - the channel message is routed to a wrong amper instance")
	}
	if BatchId == nil {
		replyToBatchId := replyMessageTimedMutexMap.StoreCompute(*RepliesToMessageId, func(value *structs.TimedBatchId) *structs.TimedBatchId {
			if value == nil {
				messagesBytes := []byte("[]")
				size := int64(len(messagesBytes))
				date := datetime.FormatDate(time.Now())
				success, fileMetadata, errU := Upload(userID, nil, &messagesBytes, util.PointerString(date), util.PointerString("application/json"), &size, util.PointerString(filepath.Join(CHAT_CHANNEL_PATH, *To)), nil)
				if success && errU == nil {
					value = &structs.TimedBatchId{
						BatchId: fileMetadata.Id,
						Time:    time.Now(),
					}
				} else {
					util.Loggify(errU)
					return nil
				}
			} else {
				value.Time = time.Now()
			}
			return value
		})
		if replyToBatchId == nil {
			return false, nil, fmt.Errorf("sendChannelReplyMessage - not able to reserve a batch file, try again later")
		} else {
			BatchId = replyToBatchId.BatchId
		}
		updateReceiveChatChannel(userID, sessionId, RepliesToMessageBatchId, RepliesToMessageId, RepliesToMessageType, util.PointerString("replyInitialisation"), util.PointerString("replyInitialisation"), BatchId, From, To)
	}
	updateReceiveChatChannel(userID, sessionId, RepliesToMessageBatchId, RepliesToMessageId, RepliesToMessageType, util.PointerString("reply"), util.PointerString("reply"), From, From, To)

	messageTime := time.Now().UnixMilli()
	message = &structs.Message{
		From:        From,
		FromUser:    business.GetUser(&fromInt64, true),
		To:          To,
		Id:          Id,
		Text:        Text,
		DateTime:    messageTime,
		BatchId:     BatchId,
		Attachments: Attachments,
	}

	newChannelRepliesCS.StoreCompute(*BatchId, func(value []structs.Message) []structs.Message {
		if value != nil {
			return append(value, *message)
		} else {
			return []structs.Message{*message}
		}
	})

	userMessageUpdate := structs.UserMessageUpdate{
		Message:        message,
		MessageType:    util.PointerString("channel"),
		UpdateType:     util.PointerString("newMessage"),
		OpperationType: util.PointerString("newMessage"),
		From:           From,
		To:             message.To,
	}
	//send update to all participants with a defer to not bloeck the current thread
	defer sendUpdateToChannelParticipants(userID, sessionId, channel, &userMessageUpdate, "chatReply")
	return true, message, nil
}

var batchReplyMessagesTimedCache = csmap.Create[string, *structs.TimedMessages](
	// set the number of map shards. the default value is 32.
	csmap.WithShardCount[string, *structs.TimedMessages](100),

	// if don't set custom hasher, use the built-in maphash.
	csmap.WithCustomHasher[string, *structs.TimedMessages](func(key string) uint64 {
		hash := fnv.New64a()
		hash.Write([]byte(key))
		return hash.Sum64()
	}),

	// set the total capacity, every shard map has total capacity/shard count capacity. the default value is 0.
	csmap.WithSize[string, *structs.TimedMessages](10000))

// this mutex map ensures there is no concurrent access for a single batch read and write opperation
// each mutex belongs to specific batch id
var batchReplyUpdateTimedMutexMap = csmap.Create[string, *structs.TimedMutex](
	// set the number of map shards. the default value is 32.
	csmap.WithShardCount[string, *structs.TimedMutex](100),

	// if don't set custom hasher, use the built-in maphash.
	csmap.WithCustomHasher[string, *structs.TimedMutex](func(key string) uint64 {
		hash := fnv.New64a()
		hash.Write([]byte(key))
		return hash.Sum64()
	}),

	// set the total capacity, every shard map has total capacity/shard count capacity. the default value is 0.
	csmap.WithSize[string, *structs.TimedMutex](10000))

func FetchChatRepliesChannel(userID *int64, BatchId *string, MessageType *string) (result []structs.Message, participants map[int64]*structs.User, err error) {
	result = make([]structs.Message, 0)
	// this mutex map ensures there is no concurrent access for a single batch read and write opperation
	// each mutex belongs to specific batch id
	batchReplyMutex := batchReplyUpdateTimedMutexMap.StoreCompute(*BatchId, func(value *structs.TimedMutex) *structs.TimedMutex {
		if value == nil {
			value = &structs.TimedMutex{
				Mutex: &sync.RWMutex{},
				Time:  time.Now(),
			}
		} else {
			value.Time = time.Now()
		}
		return value
	})
	batchReplyMutex.Mutex.RLock()
	defer batchReplyMutex.Mutex.RUnlock()

	timedMessages := batchReplyMessagesTimedCache.StoreCompute(*BatchId, func(value *structs.TimedMessages) *structs.TimedMessages {
		if value == nil {
			reader, errR := GetFileBody(userID, BatchId)
			if errR != nil {
				util.Loggify(errR)
				return nil
			}
			existingMessagesData, errPD := io.ReadAll(*reader)
			if errPD != nil {
				util.Loggify(errPD)
				return nil
			}
			var existingMessages []structs.Message
			errUM := json.Unmarshal(existingMessagesData, &existingMessages)
			if errUM != nil {
				util.Loggify(errUM)
				return nil
			}

			value = &structs.TimedMessages{
				Time:     time.Now(),
				Messages: &existingMessages,
			}
		} else {
			value.Time = time.Now()
		}
		return value
	})
	var intermediateResult []structs.Message = make([]structs.Message, 0)
	for _, existingMessage := range *timedMessages.Messages {
		if !existingMessage.Deleted {
			intermediateResult = append(intermediateResult, existingMessage)
		}
	}

	newRepliesToAppend, okNR := newChannelRepliesCS.Load(*BatchId)
	if okNR {
		intermediateResult = append(intermediateResult, newRepliesToAppend...)
	}
	participants = make(map[int64]*structs.User)
	channelRepliesMessageUpdatesCS.LoadLocked(*BatchId, func(value map[string][]structs.MessageUpdate) {
	FirstLoop:
		for i := 0; i < len(intermediateResult); i++ {
			messageUpdates := value[*intermediateResult[i].Id]
			for _, messageUpdate := range messageUpdates {
				fromInt64, errFI := strconv.ParseInt(*messageUpdate.From, 10, 64)
				if errFI != nil {
					util.Loggify(errFI)
					continue
				}
				if *messageUpdate.UpdateType == "reaction" {
					if intermediateResult[i].Reactions == nil {
						intermediateResult[i].Reactions = map[string][]int64{}
					}
					if *messageUpdate.OpperationType == "add" {
						if intermediateResult[i].Reactions[*messageUpdate.Value] == nil {
							intermediateResult[i].Reactions[*messageUpdate.Value] = make([]int64, 0)
							intermediateResult[i].Reactions[*messageUpdate.Value] = append(intermediateResult[i].Reactions[*messageUpdate.Value], fromInt64)
						} else {
							if !arrays.Contains(intermediateResult[i].Reactions[*messageUpdate.Value], fromInt64) {
								intermediateResult[i].Reactions[*messageUpdate.Value] = append(intermediateResult[i].Reactions[*messageUpdate.Value], fromInt64)
							}
						}
					} else if *messageUpdate.OpperationType == "remove" {
						if intermediateResult[i].Reactions[*messageUpdate.Value] != nil {
							intermediateResult[i].Reactions[*messageUpdate.Value] = arrays.Remove(intermediateResult[i].Reactions[*messageUpdate.Value], fromInt64)
						}
					}
				} else if *messageUpdate.UpdateType == "edit" {
					intermediateResult[i].Text = messageUpdate.Value
				} else if *messageUpdate.UpdateType == "remove" {
					continue FirstLoop
				} else if *messageUpdate.UpdateType == "reply" {
					if intermediateResult[i].Replies == nil {
						intermediateResult[i].Replies = make(map[string]int)
					}
					intermediateResult[i].Replies[*messageUpdate.Value]++
				} else if *messageUpdate.UpdateType == "replyInitialisation" {
					intermediateResult[i].ReplyBatchId = messageUpdate.Value
				}
			}
			//add the from user if missing in the participants map
			if intermediateResult[i].From != nil {
				fromUserInt64, _ := strconv.ParseInt(*intermediateResult[i].From, 10, 64)
				if participants[fromUserInt64] == nil {
					participants[fromUserInt64] = business.GetUser(&fromUserInt64, true)
				}
			}
			//add users reacting to the message to the participants map
			if intermediateResult[i].Reactions != nil {
				for _, userIds := range intermediateResult[i].Reactions {
					if len(userIds) > 0 {
						for _, userId := range userIds {
							if participants[userId] == nil {
								participants[userId] = business.GetUser(&userId, true)
							}
						}
					}
				}
			}
			result = append(result, intermediateResult[i])
		}
	})

	return result, participants, nil
}

var channelRepliesMessageUpdatesCS = csmap.Create[string, map[string][]structs.MessageUpdate](
	// set the number of map shards. the default value is 32.
	csmap.WithShardCount[string, map[string][]structs.MessageUpdate](100),

	// if don't set custom hasher, use the built-in maphash.
	csmap.WithCustomHasher[string, map[string][]structs.MessageUpdate](func(key string) uint64 {
		hash := fnv.New64a()
		hash.Write([]byte(key))
		return hash.Sum64()
	}),

	// set the total capacity, every shard map has total capacity/shard count capacity. the default value is 0.
	csmap.WithSize[string, map[string][]structs.MessageUpdate](10000))

func UpdateChatReplyChannel(userID *int64, sessionId *string, BatchId *string, MessageId *string, MessageType *string, UpdateType *string, OpperationType *string, Value *string, From *string, To *string, Participants *[]string) (success bool, err error) {
	channelId, errTI := strconv.ParseInt(*To, 10, 64)
	if errTI != nil {
		return false, fmt.Errorf("UpdateChatReplyChannel - the 'to' parameter supplied doesn't resolve to a valid identifying number")
	}
	fromInt64, errFI := strconv.ParseInt(*From, 10, 64)
	if errFI != nil {
		return false, fmt.Errorf("UpdateChatReplyChannel - the 'from' parameter supplied doesn't resolve to a valid identifying number")
	}
	channel := business.GetChatChannel(&channelId)
	if *channel.AmperId != *business.AmperId() {
		return false, fmt.Errorf("UpdateChatReplyChannel - the channel message is routed to a wrong amper instance")
	}
	messageUpdate := structs.MessageUpdate{
		MessageType:    MessageType,
		UpdateType:     UpdateType,
		OpperationType: OpperationType,
		From:           From,
		To:             To,
		Value:          Value,
	}

	channelRepliesMessageUpdatesCS.StoreCompute(*BatchId, func(value map[string][]structs.MessageUpdate) map[string][]structs.MessageUpdate {
		if value != nil {
			if value[*MessageId] != nil {
				value[*MessageId] = append(value[*MessageId], messageUpdate)
			} else {
				value[*MessageId] = []structs.MessageUpdate{messageUpdate}
			}
		} else {
			value = make(map[string][]structs.MessageUpdate)
			value[*MessageId] = []structs.MessageUpdate{messageUpdate}
		}
		return value
	})

	userMessageUpdate := structs.UserMessageUpdate{
		Message: &structs.Message{
			Id:       MessageId,
			BatchId:  BatchId,
			FromUser: business.GetUser(&fromInt64, true),
		},
		MessageType:    MessageType,
		UpdateType:     UpdateType,
		OpperationType: OpperationType,
		From:           From,
		To:             To,
		Value:          Value,
	}
	//send update to all participants with a defer to not bloeck the current thread
	defer sendUpdateToChannelParticipants(userID, sessionId, channel, &userMessageUpdate, "chatReply")
	return true, nil
}
