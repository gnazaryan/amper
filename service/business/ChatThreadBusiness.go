package business

import (
	"amper/cache/business"
	"amper/common/structs"
	"amper/common/util"
	"amper/common/util/arrays"
	"amper/common/util/datetime"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	csmap "github.com/gnazaryan/concurrent-swiss-map"
)

var threadMutexTimoutDuration time.Duration = 30 * time.Minute

func ClearThreadMutexMaps() {
	timedoutKeys := []string{}
	threadBatchUpdateTimedMutexMap.Range(func(key string, value *structs.TimedMutex) (stop bool) {
		if value != nil {
			if time.Now().After(value.Time.Add(threadMutexTimoutDuration)) {
				timedoutKeys = append(timedoutKeys, key)
			}
		}
		return false
	})
	for _, key := range timedoutKeys {
		threadBatchUpdateTimedMutexMap.Delete(key)
	}
}

var threadReplyMutexTimoutDuration time.Duration = 30 * time.Minute

func ClearThreadReplyMutexMaps() {
	timedoutKeys := []string{}
	threadBatchReplyUpdateTimedMutexMap.Range(func(key string, value *structs.TimedMutex) (stop bool) {
		if value != nil {
			if time.Now().After(value.Time.Add(threadReplyMutexTimoutDuration)) {
				timedoutKeys = append(timedoutKeys, key)
			}
		}
		return false
	})
	for _, key := range timedoutKeys {
		threadBatchReplyUpdateTimedMutexMap.Delete(key)
	}
}

var newBatchIdMutexTimoutDuration time.Duration = 30 * time.Minute

func ClearThreadNewBatchIdMutexMaps() {
	timedoutKeys := []string{}
	threadNewBatchTimedMutexMap.Range(func(key string, value *structs.TimedMutex) (stop bool) {
		if value != nil {
			if time.Now().After(value.Time.Add(newBatchIdMutexTimoutDuration)) {
				timedoutKeys = append(timedoutKeys, key)
			}
		}
		return false
	})
	for _, key := range timedoutKeys {
		threadNewBatchTimedMutexMap.Delete(key)
	}
}

var cacheThreadMessagesTimoutDuration time.Duration = 30 * time.Minute

func ClearThreadTimedoutMessages() {
	timedoutKeys := []string{}
	batchMessagesTimedCache.Range(func(key string, value *structs.TimedMessages) (stop bool) {
		if value != nil {
			if time.Now().After(value.Time.Add(cacheThreadMessagesTimoutDuration)) {
				timedoutKeys = append(timedoutKeys, key)
			}
		}
		return false
	})
	for _, key := range timedoutKeys {
		batchMessagesTimedCache.Delete(key)
	}
}

var cacheThreadReplyMessagesTimoutDuration time.Duration = 30 * time.Minute

func ClearThreadReplyTimedoutMessages() {
	timedoutKeys := []string{}
	threadBatchMessagesReplyTimedCache.Range(func(key string, value *structs.TimedMessages) (stop bool) {
		if value != nil {
			if time.Now().After(value.Time.Add(cacheThreadReplyMessagesTimoutDuration)) {
				timedoutKeys = append(timedoutKeys, key)
			}
		}
		return false
	})
	for _, key := range timedoutKeys {
		threadBatchMessagesReplyTimedCache.Delete(key)
	}
}

var newThreadMessageCS = csmap.Create[string, []structs.Message](
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

func SaveThreadHistory() {
	threadChatDirectory, errTHD := GetThreadChatDirectory()
	if errTHD != nil {
		util.Loggify(errTHD)
		return
	}
	newThreadMessageCS.RangeDelete(func(key string, value []structs.Message) (stop bool) {
		if len(value) > 0 {
			//this mutex map ensures there is no concurrent access for a single batch read and write opperation
			//each mutex belongs to specific batch id
			batchMutex := threadBatchUpdateTimedMutexMap.StoreCompute(key, func(value *structs.TimedMutex) *structs.TimedMutex {
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
			existingMessages = append(existingMessages, value...)

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
			//to the history file
			if len(updatedMessages) > CHAT_FILE_CHUNK_SIZE {
				//All messages will have the To value as the thread id
				//and the To value maps to the thread mutex
				//take a RW lock to that specific mutex to update the history file
				to := value[0].To
				//befor initializing a new batch file take a lock so the history file reading thread
				//reads the lates batch id
				batchMutex := threadNewBatchTimedMutexMap.StoreCompute(*to, func(value *structs.TimedMutex) *structs.TimedMutex {
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
				historyFile := filepath.Join(*threadChatDirectory, *to, "history")
				historyFileBytes, errHFB := os.ReadFile(historyFile)
				if errHFB != nil {
					util.Loggify(errHFB)
					return false
				}
				history := &structs.ChatHistory{}
				errH := json.Unmarshal(historyFileBytes, history)
				if errH != nil {
					util.Loggify(errH)
					return false
				}
				messagesBytes := []byte("[]")
				size := int64(len(messagesBytes))
				date := datetime.FormatDate(time.Now())
				success, fileMetadata, errU := Upload(&batchIdStruct.UserId, nil, &messagesBytes, util.PointerString(date), util.PointerString("application/json"), &size, util.PointerString(filepath.Join(CHAT_PATH, *to)), nil)
				if success && errU == nil {
					//Mark the last history item to be full, since it exceeded the max file size
					history.HistoryItems[len(history.HistoryItems)-1].Full = true
					history.HistoryItems = append(history.HistoryItems, structs.ChatHistoryItem{
						Id:   fileMetadata.Id,
						Full: false,
					})

					historyBytes, errHB := json.Marshal(history)
					if errHB != nil {
						util.Loggify(errHB)
						return false
					}

					errR := os.Remove(historyFile)
					if errR != nil {
						util.Loggify(errR)
						return false
					}

					errH := os.WriteFile(historyFile, historyBytes, 0644)
					if errH != nil {
						util.Loggify(errH)
						return false
					}
					instance := business.GetAmperInstance(*business.AmperId())
					sessionId := GenerateSessionId(instance.Identifier, instance.Key, util.PointerString("app"))
					//Send update to all active participants for the new batch added
					//All active ui users should add the batch id to the chat item in client side
					userMessageUpdate := structs.UserMessageUpdate{
						MessageType:    util.PointerString("thread"),
						UpdateType:     util.PointerString("newBatch"),
						OpperationType: util.PointerString("newBatch"),
						From:           util.PointerString(strconv.FormatInt(SYSTEM_USER_ID, 10)),
						To:             to,
						Value:          fileMetadata.Id,
					}
					//send update to all paricipants for the new batch added
					defer sendUpdateToAllParticipants(&SYSTEM_USER_ID, &sessionId, util.PointerString("chat"), history, &userMessageUpdate)
				}
			}
		}
		return false
	})
}

func SaveThreadMessageUpdates() {
	threadMessageUpdatesCS.RangeDelete(func(key string, value map[string][]structs.MessageUpdate) (stop bool) {
		//this mutex map ensures there is no concurrent access for a single batch read and write opperation
		//each mutex belongs to specific batch id
		batchMutex := threadBatchUpdateTimedMutexMap.StoreCompute(key, func(value *structs.TimedMutex) *structs.TimedMutex {
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

func SaveThreadRepliesHistory() {
	newThreadRepliesCS.RangeDelete(func(key string, value []structs.Message) (stop bool) {
		//this mutex map ensures there is no concurrent access for a single batch read and write opperation
		//each mutex belongs to specific batch id
		batchReplyMutex := threadBatchReplyUpdateTimedMutexMap.StoreCompute(key, func(value *structs.TimedMutex) *structs.TimedMutex {
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
		threadBatchMessagesReplyTimedCache.Delete(key)

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

func SaveThreadRepliesMessageUpdates() {
	threadRepliesMessageUpdatesCS.RangeDelete(func(key string, value map[string][]structs.MessageUpdate) (stop bool) {
		//this mutex map ensures there is no concurrent access for a single batch read and write opperation
		//each mutex belongs to specific batch id
		batchReplyMutex := threadBatchReplyUpdateTimedMutexMap.StoreCompute(key, func(value *structs.TimedMutex) *structs.TimedMutex {
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

		threadBatchMessagesReplyTimedCache.Delete(key)

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

// this mutex map ensures there is no concurrent access for a batch id while there is a process of initiation  of a new batch id
// each mutex belongs to specific batch id
var threadNewBatchTimedMutexMap = csmap.Create[string, *structs.TimedMutex](
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

func sendThreadMessage(userID *int64, sessionId *string, from *string, to *string, participants *[]int64, Type *string, id *string, text *string, Attachments *[]structs.ChatAttachment) (successResult bool, message *structs.Message, errResult error) {
	fromInt64, errFI := strconv.ParseInt(*from, 10, 64)
	if errFI != nil {
		return false, nil, fmt.Errorf("the from parameter supplied doesn't resolve to a valid identifying number")
	}
	if *userID != fromInt64 {
		return false, nil, fmt.Errorf("sending user must be the one currently logged in to the system")
	}
	threadChatDirectory, errTHD := GetThreadChatDirectory()
	if errTHD != nil {
		util.Loggify(errTHD)
		return false, nil, fmt.Errorf("not able to retrieve the thread directory")
	}
	var BatchId *string
	historyFile := filepath.Join(*threadChatDirectory, *to, "history")
	messageTime := time.Now().UnixMilli()
	if _, err := os.Stat(historyFile); errors.Is(err, os.ErrNotExist) {
		messagesBytes := []byte("[]")
		size := int64(len(messagesBytes))
		date := datetime.FormatDate(time.Now())
		success, fileMetadata, errU := Upload(userID, nil, &messagesBytes, util.PointerString(date), util.PointerString("application/json"), &size, util.PointerString(filepath.Join(CHAT_PATH, *to)), nil)
		if success && errU == nil {
			BatchId = fileMetadata.Id
			history := structs.ChatHistory{
				From: from,
				To:   to,
				HistoryItems: []structs.ChatHistoryItem{{
					Id:   fileMetadata.Id,
					Full: false,
				}},
				LastUpdateTime: time.Now().UnixMilli(),
				Participants:   participants,
				UnreadMessages: 1,
			}
			historyBytes, errHB := json.Marshal(history)
			if errHB != nil {
				util.Loggify(errHB)
				return false, nil, fmt.Errorf("sendThreadMessage - not able to initilize the history file due to json error")
			}

			errChD := os.MkdirAll(filepath.Join(*threadChatDirectory, *to), os.ModePerm)
			if errChD != nil {
				util.Loggify(errChD)
				return false, nil, fmt.Errorf("sendThreadMessage - not able to initilize the history file directory")
			}

			errH := os.WriteFile(historyFile, historyBytes, 0644)
			if errH != nil {
				util.Loggify(errH)
				return false, nil, fmt.Errorf("sendThreadMessage - not able to initilize the history file")
			}
			message = &structs.Message{
				From:        from,
				To:          to,
				Id:          id,
				Text:        text,
				DateTime:    messageTime,
				BatchId:     BatchId,
				Attachments: Attachments,
			}
			_, errITC := initThreadChatForParticipants(userID, sessionId, to, participants, &history, message)
			if errITC != nil {
				util.Loggify(errITC)
				return false, nil, fmt.Errorf("not able to initilize the chat thread for all participants")
			}
		} else {
			util.Loggify(errU)
			return false, nil, fmt.Errorf("sendThreadMessage - not able to initilize the history file")
		}
	} else {
		//Take a read lock for the current thread history file read operation
		//because a chat history save thread may be working on writing it at the time
		//this code runs
		batchMutex := threadNewBatchTimedMutexMap.StoreCompute(*to, func(value *structs.TimedMutex) *structs.TimedMutex {
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

		historyFileBytes, errHFB := os.ReadFile(historyFile)
		if errHFB != nil {
			util.Loggify(errHFB)
			return false, nil, fmt.Errorf("sendThreadMessage - not able to read the history file")
		}
		history := &structs.ChatHistory{}
		errH := json.Unmarshal(historyFileBytes, history)
		if errH != nil {
			util.Loggify(errH)
			return false, nil, fmt.Errorf("sendThreadMessage - not able to parse the history file")
		}
		BatchId = history.HistoryItems[len(history.HistoryItems)-1].Id
		message = &structs.Message{
			From:        from,
			To:          to,
			Id:          id,
			Text:        text,
			DateTime:    messageTime,
			BatchId:     BatchId,
			Attachments: Attachments,
		}
		_, errITC := updateThreadChatForParticipants(userID, sessionId, to, participants, history, message)
		if errITC != nil {
			util.Loggify(errITC)
			return false, nil, fmt.Errorf("not able to update the chat thread for all participants")
		}
	}

	newThreadMessageCS.StoreCompute(*BatchId, func(value []structs.Message) []structs.Message {
		if value != nil {
			return append(value, *message)
		} else {
			return []structs.Message{*message}
		}
	})
	return true, message, nil
}

func updateThreadChatForParticipants(userID *int64, sessionId *string, threadId *string, participants *[]int64, ChatHistory *structs.ChatHistory, message *structs.Message) (success bool, err error) {
	success = true
	var instanceToParticipantsMap map[int64][]int64 = make(map[int64][]int64)
	for _, participant := range *participants {
		if participant == *userID {
			//skip providing update for the sending user, since sending user has the update
			continue
		}
		user := business.GetUser(&participant, true)
		if instanceToParticipantsMap[*user.AmperId] == nil {
			instanceToParticipantsMap[*user.AmperId] = []int64{*user.ID}
		} else {
			instanceToParticipantsMap[*user.AmperId] = append(instanceToParticipantsMap[*user.AmperId], *user.ID)
		}
	}

	for instanceId, instanceParticipants := range instanceToParticipantsMap {
		instance := business.GetAmperInstance(instanceId)
		successInstance, _, errPS := DedicatedCallWithRetry(userID, sessionId, map[string]string{
			"amperInstance": "chat/receiveThread",
		}, map[string]interface{}{
			"threadId":     threadId,
			"participants": instanceParticipants,
			"chatHistory":  ChatHistory,
			"message":      message,
		}, instance)
		if errPS != nil || !successInstance {
			util.Loggify(errPS)
			success = false
		}
	}
	return success, nil
}

func initThreadChatForParticipants(userID *int64, sessionId *string, threadId *string, participants *[]int64, ChatHistory *structs.ChatHistory, message *structs.Message) (success bool, err error) {
	success = true
	var instanceToParticipantsMap map[int64][]int64 = make(map[int64][]int64)
	label := ""
	for _, participant := range *participants {
		user := business.GetUser(&participant, true)
		if len(label) > 0 {
			label += ", "
		}
		label += *user.FirstName
		if instanceToParticipantsMap[*user.AmperId] == nil {
			instanceToParticipantsMap[*user.AmperId] = []int64{*user.ID}
		} else {
			instanceToParticipantsMap[*user.AmperId] = append(instanceToParticipantsMap[*user.AmperId], *user.ID)
		}
	}

	if len(label) > 25 {
		label = label[0:25] + "..."
	}
	for instanceId, instanceParticipants := range instanceToParticipantsMap {
		instance := business.GetAmperInstance(instanceId)
		successInstance, _, errPS := DedicatedCallWithRetry(userID, sessionId, map[string]string{
			"amperInstance": "chat/initThread",
		}, map[string]interface{}{
			"label":                label,
			"threadId":             threadId,
			"participants":         participants,
			"instanceParticipants": instanceParticipants,
			"chatHistory":          ChatHistory,
			"message":              message,
		}, instance)
		if errPS != nil || !successInstance {
			util.Loggify(errPS)
			success = false
		}
	}

	return success, err
}

func InitChatThread(userID *int64, Label *string, ThreadId *string, Participants *[]int64, InstanceParticipants *[]int64, ChatHistory *structs.ChatHistory, Message *structs.Message) (success bool, err error) {
	participants := []structs.User{}
	if ChatHistory.Participants != nil {
		for _, participant := range *ChatHistory.Participants {
			participants = append(participants, *business.GetUser(&participant, true))
		}
	}
	for _, instanceParticipant := range *InstanceParticipants {
		chatDirectory, errCh := GetChatDirectory(&instanceParticipant)
		if errCh != nil {
			util.Loggify(errCh)
			continue
		}
		historyFile := filepath.Join(*chatDirectory, "threads", *ThreadId, "history")
		errRF := os.Remove(historyFile)
		if errRF != nil && !errors.Is(errRF, os.ErrNotExist) {
			util.Loggify(errRF)
			continue
		}
		chatThreadHistory := structs.ChatThreadHistory{
			Label:          Label,
			ThreadId:       ThreadId,
			UnreadMessages: 1,
			LastUpdateTime: time.Now().UnixMilli(),
		}

		historyBytes, errHB := json.Marshal(chatThreadHistory)
		if errHB != nil {
			util.Loggify(errHB)
			continue
		}

		errChD := os.MkdirAll(filepath.Join(*chatDirectory, "threads", *ThreadId), os.ModePerm)
		if errChD != nil {
			util.Loggify(errChD)
			continue
		}

		errH := os.WriteFile(historyFile, historyBytes, 0644)
		if errH != nil {
			util.Loggify(errH)
			continue
		}

		var userUpdate interface{} = structs.UserMessageUpdate{
			Message:        Message,
			MessageType:    util.PointerString("thread"),
			UpdateType:     util.PointerString("newMessage"),
			OpperationType: util.PointerString("newMessage"),
			From:           Message.From,
			To:             Message.To,
			ChatHistory:    ChatHistory,
			Users:          &participants,
		}
		PutUpdate(&instanceParticipant, util.PointerString("chat"), &userUpdate)
	}
	return true, nil
}

func ReceiveThread(userID *int64, ThreadId *string, Participants *[]int64, ChatHistory *structs.ChatHistory, Message *structs.Message) (success bool, err error) {
	participants := []structs.User{}
	if ChatHistory.Participants != nil {
		for _, participant := range *ChatHistory.Participants {
			participants = append(participants, *business.GetUser(&participant, true))
		}
	}
	for _, participant := range *Participants {
		chatDirectory, errCh := GetChatDirectory(&participant)
		if errCh != nil {
			util.Loggify(errCh)
			continue
		}
		historyFile := filepath.Join(*chatDirectory, "threads", *ThreadId, "history")

		historyFileBytes, errHFB := os.ReadFile(historyFile)
		if errHFB != nil {
			util.Loggify(errHFB)
			continue
		}
		chatThreadHistory := &structs.ChatThreadHistory{}
		errH := json.Unmarshal(historyFileBytes, chatThreadHistory)
		if errH != nil {
			util.Loggify(errH)
			continue
		}
		chatThreadHistory.UnreadMessages++
		chatThreadHistory.LastUpdateTime = time.Now().UnixMilli()

		historyBytes, errHB := json.Marshal(chatThreadHistory)
		if errHB != nil {
			util.Loggify(errHB)
			continue
		}

		errChD := os.MkdirAll(filepath.Join(*chatDirectory, "threads", *ThreadId), os.ModePerm)
		if errChD != nil {
			util.Loggify(errChD)
			continue
		}

		errH = os.WriteFile(historyFile, historyBytes, 0644)
		if errH != nil {
			util.Loggify(errH)
			continue
		}
		ChatHistory.UnreadMessages = chatThreadHistory.UnreadMessages
		var userUpdate interface{} = structs.UserMessageUpdate{
			Message:        Message,
			MessageType:    util.PointerString("thread"),
			UpdateType:     util.PointerString("newMessage"),
			OpperationType: util.PointerString("newMessage"),
			From:           Message.From,
			To:             Message.To,
			ChatHistory:    ChatHistory,
			Users:          &participants,
		}
		PutUpdate(&participant, util.PointerString("chat"), &userUpdate)
	}
	return true, nil
}

func GetThreadChatDirectory() (result *string, err error) {
	instance := business.GetAmperInstance(*business.AmperId())
	rootDriveDirectory := instance.Directory
	if rootDriveDirectory == nil || len(*rootDriveDirectory) < 1 {
		return result, fmt.Errorf("empty directory found, amper instance %d is not configured for a directory", *business.AmperId())
	}
	result = util.PointerString(filepath.Join(*rootDriveDirectory, "chat", "thread"))
	errUD := os.MkdirAll(*result, os.ModePerm)
	if errUD != nil && !errors.Is(err, os.ErrExist) {
		util.Loggify(errUD)
		return result, fmt.Errorf("not able to locate the user's active directory in chat '%s', please contect the support", *result)
	}
	return result, nil
}

func GetThreadsHistories(userID *int64, ThreadIds *[]string) (result *[]structs.ChatHistory, err error) {
	chatHistories := make([]structs.ChatHistory, 0)
	threadChatDirectory, errTHD := GetThreadChatDirectory()
	if errTHD != nil {
		util.Loggify(errTHD)
		return nil, fmt.Errorf("GetThreadsHistories - not able to retrieve the thread directory")
	}

	if ThreadIds != nil {
		for _, threadId := range *ThreadIds {
			historyFile := filepath.Join(*threadChatDirectory, threadId, "history")
			if _, err := os.Stat(historyFile); !errors.Is(err, os.ErrNotExist) {
				historyFileBytes, errHFB := os.ReadFile(historyFile)
				if errHFB != nil {
					util.Loggify(errHFB)
					continue
				}
				history := &structs.ChatHistory{}
				errH := json.Unmarshal(historyFileBytes, history)
				if errH != nil {
					util.Loggify(errH)
					continue
				}
				chatHistories = append(chatHistories, *history)
			}
		}
	}
	result = &chatHistories
	return result, nil
}

// this mutex map ensures there is no concurrent access for a single batch read and write opperation
// each mutex belongs to specific batch id
var threadBatchUpdateTimedMutexMap = csmap.Create[string, *structs.TimedMutex](
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

var threadBatchMessagesTimedCache = csmap.Create[string, *structs.TimedMessages](
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

func fetchChatThread(userID *int64, from *string, to *string, id *string, includeLatest bool) (result []structs.Message, participants map[int64]*structs.User, err error) {
	result = make([]structs.Message, 0)
	// this mutex map ensures there is no concurrent access for a single batch read and write opperation
	// each mutex belongs to specific batch id
	batchMutex := threadBatchUpdateTimedMutexMap.StoreCompute(*id, func(value *structs.TimedMutex) *structs.TimedMutex {
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
		newMessages, ok := newThreadMessageCS.Load(*id)
		if ok {
			intermediateResult = append(intermediateResult, newMessages...)
		}
	}

	threadMessageUpdatesCS.LoadLocked(*id, func(value map[string][]structs.MessageUpdate) {
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
			result = append(result, intermediateResult[i])
		}
	})

	return result, nil, nil
}

func markChatUnreadThread(userID *int64, to *string, chatType *string) (success bool, err error) {
	chatDirectory, errCh := GetChatDirectory(userID)
	if errCh != nil {
		return false, fmt.Errorf("markChatUnreadThread - not able to locate chat directory")
	}
	historyFile := filepath.Join(*chatDirectory, "threads", *to, "history")
	if _, err := os.Stat(historyFile); !errors.Is(err, os.ErrNotExist) {
		historyFileBytes, errHFB := os.ReadFile(historyFile)
		if errHFB != nil {
			util.Loggify(errHFB)
			return false, fmt.Errorf("markChatUnreadThread - not able to read the history file")
		}
		history := &structs.ChatThreadHistory{}
		errH := json.Unmarshal(historyFileBytes, history)
		if errH != nil {
			util.Loggify(errH)
			return false, fmt.Errorf("markChatUnreadThread - not able to parse the history file")
		}
		history.UnreadMessages = 0
		historyBytes, errHB := json.Marshal(history)
		if errHB != nil {
			util.Loggify(errHB)
			return false, fmt.Errorf("markChatUnreadThread - not able to initilize the history file due to json error 1")
		}

		errR := os.Remove(historyFile)
		if errR != nil {
			util.Loggify(errR)
			return false, fmt.Errorf("markChatUnreadThread - not able to remove the history file for rewriting")
		}
		errChD := os.MkdirAll(filepath.Join(*chatDirectory, "directs", *to), os.ModePerm)
		if errChD != nil {
			util.Loggify(errChD)
			return false, fmt.Errorf("markChatUnreadThread - not able to initilize the history file directory 1")
		}

		errH = os.WriteFile(historyFile, historyBytes, 0644)
		if errH != nil {
			util.Loggify(errH)
			return false, fmt.Errorf("markChatUnreadThread - not able to initilize the history file 1")
		}
	}
	return true, nil
}

var threadMessageUpdatesCS = csmap.Create[string, map[string][]structs.MessageUpdate](
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

func updateReceiveChatThread(userID *int64, sessionId *string, BatchId *string, MessageId *string, MessageType *string, UpdateType *string, OpperationType *string, Value *string, From *string, To *string) (success bool, err error) {
	messageUpdate := structs.MessageUpdate{
		MessageType:    MessageType,
		UpdateType:     UpdateType,
		OpperationType: OpperationType,
		From:           From,
		To:             To,
		Value:          Value,
	}

	threadMessageUpdatesCS.StoreCompute(*BatchId, func(value map[string][]structs.MessageUpdate) map[string][]structs.MessageUpdate {
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

	threadChatDirectory, errTHD := GetThreadChatDirectory()
	if errTHD != nil {
		util.Loggify(errTHD)
		return false, fmt.Errorf("updateChatThread - not able to retrieve the thread directory")
	}
	historyFile := filepath.Join(*threadChatDirectory, *To, "history")
	if _, err := os.Stat(historyFile); !errors.Is(err, os.ErrNotExist) {
		historyFileBytes, errHFB := os.ReadFile(historyFile)
		if errHFB != nil {
			util.Loggify(errHFB)
			return false, fmt.Errorf("updateChatThread - not able to read the history file")
		}
		history := &structs.ChatHistory{}
		errH := json.Unmarshal(historyFileBytes, history)
		if errH != nil {
			util.Loggify(errH)
			return false, fmt.Errorf("updateChatThread - not able to parse the history file")
		}
		userMessageUpdate := structs.UserMessageUpdate{
			Message: &structs.Message{
				Id: MessageId,
			},
			MessageType:    MessageType,
			UpdateType:     UpdateType,
			OpperationType: OpperationType,
			From:           From,
			To:             To,
			Value:          Value,
		}
		sendUpdateToAllParticipants(userID, sessionId, util.PointerString("chat"), history, &userMessageUpdate)
	} else {
		return false, fmt.Errorf("updateChatThread - not able to locate the history file for the thread")
	}
	return true, nil
}

func sendUpdateToAllParticipants(userID *int64, sessionId *string, category *string, history *structs.ChatHistory, messageUpdate *structs.UserMessageUpdate) {
	var instanceToParticipantsMap map[int64][]int64 = make(map[int64][]int64)
	for _, participant := range *history.Participants {
		if *userID == participant {
			//skipp the updating user, since it already has the update
			continue
		}
		user := business.GetUser(&participant, true)
		if instanceToParticipantsMap[*user.AmperId] == nil {
			instanceToParticipantsMap[*user.AmperId] = []int64{*user.ID}
		} else {
			instanceToParticipantsMap[*user.AmperId] = append(instanceToParticipantsMap[*user.AmperId], *user.ID)
		}
	}

	for instanceId, instanceParticipants := range instanceToParticipantsMap {
		instance := business.GetAmperInstance(instanceId)
		_, _, errPS := DedicatedCallWithRetry(userID, sessionId, map[string]string{
			"amperInstance": "chat/updateThread",
		}, map[string]interface{}{
			"category":      category,
			"threadId":      history.To,
			"participants":  instanceParticipants,
			"messageUpdate": messageUpdate,
		}, instance)
		if errPS != nil {
			util.Loggify(errPS)
		}
	}
}

func UpdateThread(userID *int64, Category *string, ThreadId *string, Participants *[]int64, MessageUpdate *structs.UserMessageUpdate) (success bool, err error) {
	if Participants != nil {
		var messageUpdateInterface interface{} = *MessageUpdate
		for _, participant := range *Participants {
			PutUpdate(&participant, Category, &messageUpdateInterface)
		}
	}
	return true, nil
}

var newThreadRepliesCS = csmap.Create[string, []structs.Message](
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

func sendThreadReplyMessage(userID *int64, sessionId *string, BatchId *string, From *string, To *string, RepliesToMessageId *string, RepliesToMessageBatchId *string, RepliesToMessageBatchId1 *string, RepliesToMessageType *string, Id *string, Text *string, Attachments *[]structs.ChatAttachment, Participants *[]string) (success bool, message *structs.Message, err error) {
	if BatchId == nil {
		messagesBytes := []byte("[]")
		size := int64(len(messagesBytes))
		date := datetime.FormatDate(time.Now())
		success, fileMetadata, errU := Upload(userID, nil, &messagesBytes, util.PointerString(date), util.PointerString("application/json"), &size, util.PointerString(filepath.Join(CHAT_PATH, *To)), nil)
		if success && errU == nil {
			BatchId = fileMetadata.Id
		} else {
			util.Loggify(errU)
			return false, nil, fmt.Errorf("not able to reserve a batch file, try again later")
		}
		success, errUC := updateChatRemote(userID, sessionId, RepliesToMessageBatchId, RepliesToMessageId, RepliesToMessageType, util.PointerString("replyInitialisation"), util.PointerString("replyInitialisation"), BatchId, From, To)
		if errUC != nil || !success {
			util.Loggify(errUC)
			return false, nil, fmt.Errorf("not able to update the reserved batch file, try again later")
		}
	}
	success, errUC := updateChatRemote(userID, sessionId, RepliesToMessageBatchId, RepliesToMessageId, RepliesToMessageType, util.PointerString("reply"), util.PointerString("reply"), From, From, To)
	if errUC != nil || !success {
		util.Loggify(errUC)
	}
	messageTime := time.Now().UnixMilli()
	message = &structs.Message{
		From:        From,
		To:          To,
		Id:          Id,
		Text:        Text,
		DateTime:    messageTime,
		BatchId:     BatchId,
		Attachments: Attachments,
	}

	newThreadRepliesCS.StoreCompute(*BatchId, func(value []structs.Message) []structs.Message {
		if value != nil {
			return append(value, *message)
		} else {
			return []structs.Message{*message}
		}
	})

	var messageUpdate = structs.UserMessageUpdate{
		Message:        message,
		MessageType:    util.PointerString("thread"),
		UpdateType:     util.PointerString("newMessage"),
		OpperationType: util.PointerString("newMessage"),
		From:           From,
		To:             To,
	}
	SendUpdateToAllParticipants(userID, sessionId, Participants, &messageUpdate)
	return true, message, nil
}

// this mutex map ensures there is no concurrent access for a single batch read and write opperation
// each mutex belongs to specific batch id
var threadBatchReplyUpdateTimedMutexMap = csmap.Create[string, *structs.TimedMutex](
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

var threadBatchMessagesReplyTimedCache = csmap.Create[string, *structs.TimedMessages](
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

func FetchChatRepliesThread(userID *int64, BatchId *string, MessageType *string) (result []structs.Message, participants map[int64]*structs.User, err error) {
	result = make([]structs.Message, 0)
	// this mutex map ensures there is no concurrent access for a single batch read and write opperation
	// each mutex belongs to specific batch id
	batchReplyMutex := threadBatchReplyUpdateTimedMutexMap.StoreCompute(*BatchId, func(value *structs.TimedMutex) *structs.TimedMutex {
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

	timedMessages := threadBatchMessagesReplyTimedCache.StoreCompute(*BatchId, func(value *structs.TimedMessages) *structs.TimedMessages {
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

	newRepliesToAppend, okNR := newThreadRepliesCS.Load(*BatchId)
	if okNR {
		intermediateResult = append(intermediateResult, newRepliesToAppend...)
	}

	threadRepliesMessageUpdatesCS.LoadLocked(*BatchId, func(value map[string][]structs.MessageUpdate) {
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
			result = append(result, intermediateResult[i])
		}
	})
	return result, participants, nil
}

var threadRepliesMessageUpdatesCS = csmap.Create[string, map[string][]structs.MessageUpdate](
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

func UpdateChatReplyThread(userID *int64, sessionId *string, BatchId *string, MessageId *string, MessageType *string, UpdateType *string, OpperationType *string, Value *string, From *string, To *string, Participants *[]string) (success bool, err error) {
	messageUpdate := structs.MessageUpdate{
		MessageType:    MessageType,
		UpdateType:     UpdateType,
		OpperationType: OpperationType,
		From:           From,
		To:             To,
		Value:          Value,
	}

	threadRepliesMessageUpdatesCS.StoreCompute(*BatchId, func(value map[string][]structs.MessageUpdate) map[string][]structs.MessageUpdate {
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

	threadChatDirectory, errTHD := GetThreadChatDirectory()
	if errTHD != nil {
		util.Loggify(errTHD)
		return false, fmt.Errorf("UpdateChatReplyThread - not able to retrieve the thread directory")
	}
	historyFile := filepath.Join(*threadChatDirectory, *To, "history")
	if _, err := os.Stat(historyFile); !errors.Is(err, os.ErrNotExist) {
		historyFileBytes, errHFB := os.ReadFile(historyFile)
		if errHFB != nil {
			util.Loggify(errHFB)
			return false, fmt.Errorf("UpdateChatReplyThread - not able to read the history file")
		}
		history := &structs.ChatHistory{}
		errH := json.Unmarshal(historyFileBytes, history)
		if errH != nil {
			util.Loggify(errH)
			return false, fmt.Errorf("UpdateChatReplyThread - not able to parse the history file")
		}
		userMessageUpdate := structs.UserMessageUpdate{
			Message: &structs.Message{
				Id:      MessageId,
				BatchId: BatchId,
			},
			MessageType:    MessageType,
			UpdateType:     UpdateType,
			OpperationType: OpperationType,
			From:           From,
			To:             To,
			Value:          Value,
		}
		sendUpdateToAllParticipants(userID, sessionId, util.PointerString("chatReply"), history, &userMessageUpdate)
	} else {
		return false, fmt.Errorf("UpdateChatReplyThread - not able to locate the history file for the thread")
	}
	return true, nil
}
