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
	"strings"
	"sync"
	"time"

	csmap "github.com/gnazaryan/concurrent-swiss-map"
)

// The below cache is designed to hold new messages which have not been yet persisted to the key value store
var newDirectMessageRWMutex sync.RWMutex
var newDirectMessages map[int64]map[string][]structs.Message = make(map[int64]map[string][]structs.Message)

var newDirectMessageMutex sync.Mutex
var newDirectMessageMutexes map[int64]*sync.Mutex = make(map[int64]*sync.Mutex)

// The below cach is designed to hold new reactions that have not been persisted to the key value store
var newUpdates map[string]map[string][]structs.MessageUpdate = make(map[string]map[string][]structs.MessageUpdate)
var batchUpdateMutexis map[string]*sync.Mutex = make(map[string]*sync.Mutex)
var batchUpdateMutex sync.Mutex
var batchUpdateRWMutex sync.RWMutex

func init() {
	go SaveChatHistoryLoop()
	go ClearTimedoutCach()
}

func SaveChatHistoryLoop() {
	time.AfterFunc(1*time.Minute, func() {
		SaveDirectChatHistory()
		SaveDirectMessageUpdates()
		SaveDirectMessageReplies()
		SaveDirectMessageReplyUpdates()

		SaveThreadHistory()
		SaveThreadMessageUpdates()
		SaveThreadRepliesHistory()
		SaveThreadRepliesMessageUpdates()

		SaveChannelHistory()
		SaveChannelMessageUpdates()
		SaveChannelRepliesHistory()
		SaveChannelRepliesMessageUpdates()

		ClearThreadMutexMaps()
		ClearThreadReplyMutexMaps()
		ClearThreadNewBatchIdMutexMaps()
		ClearChannelMutexMaps()
		ClearChannelReplyMutexMaps()
		SaveChatHistoryLoop()
	})
}

func ClearTimedoutCach() {
	time.AfterFunc(10*time.Minute, func() {
		ClearThreadTimedoutMessages()
		ClearThreadReplyTimedoutMessages()

		ClearChannelTimedoutMessages()
		ClearChannelTimedoutReplyMessages()

		ClearTimedoutCach()
	})
}

func SaveDirectMessageReplies() {
	newRepliesCS.RangeDelete(func(key string, value []structs.Message) (stop bool) {
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

func SaveDirectMessageReplyUpdates() {
	batchReplyUpdateRWMutex.RLock()
	batchIds := make([]string, 0, len(newReplyUpdates))
	for k := range newReplyUpdates {
		batchIds = append(batchIds, k)
	}
	batchReplyUpdateRWMutex.RUnlock()
	for _, batchId := range batchIds {
		//Acure a user specific mutex lock befor sending the message
		batchReplyUpdateMutex.Lock()
		batchReplyMutex := batchReplyUpdateMutexis[batchId]
		if batchReplyMutex == nil {
			batchReplyUpdateMutexis[batchId] = &sync.Mutex{}
			batchReplyMutex = batchReplyUpdateMutexis[batchId]
		}
		batchReplyUpdateMutex.Unlock()
		batchReplyMutex.Lock()
		saveMessageReplyUpdatesInternal(&batchId)
		batchReplyUpdateRWMutex.Lock()
		delete(newReplyUpdates, batchId)
		batchReplyUpdateRWMutex.Unlock()
		batchReplyMutex.Unlock()
	}
}

func saveMessageReplyUpdatesInternal(batchId *string) {
	batchIdStruct, errBI := structs.ParseId(batchId)
	if errBI != nil {
		util.Loggify(errBI)
		return
	}
	optionalValueSplit := strings.Split(*batchIdStruct.OptionalValue, "-")
	if len(optionalValueSplit) != 2 {
		util.Loggify(fmt.Errorf("not able to save message updates due to invalid optional value is set on id"))
		return
	}
	var directory *string
	if optionalValueSplit[0] == "direct" {
		directory = util.PointerString(filepath.Join(CHAT_PATH, optionalValueSplit[1]))
	} else {
		util.Loggify(fmt.Errorf("not able to save message updates due to invalid optional value message type is set on id"))
		return
	}
	reader, _, errGF := GetFile(&batchIdStruct.UserId, directory, batchId, &structs.Version{
		Major: 0,
		Minor: 0,
		Patch: 1,
	}, util.PointerBoolean(false))
	if errGF != nil {
		util.Loggify(errGF)
		return
	}
	existingMessagesData, errPD := io.ReadAll(*reader)
	if errPD != nil {
		util.Loggify(errPD)
		return
	}
	var existingMessages []structs.Message
	errUM := json.Unmarshal(existingMessagesData, &existingMessages)
	if errUM != nil {
		util.Loggify(errUM)
		return
	}
	batchReplyUpdateRWMutex.RLock()
	batchMessageIds := make([]string, 0, len(newReplyUpdates[*batchId]))
	for k := range newReplyUpdates[*batchId] {
		batchMessageIds = append(batchMessageIds, k)
	}
	batchReplyUpdateRWMutex.RUnlock()
	for i := 0; i < len(existingMessages); i++ {
		if arrays.Contains(batchMessageIds, *existingMessages[i].Id) {
			batchReplyUpdateRWMutex.RLock()
			messageUpdates := newReplyUpdates[*batchId][*existingMessages[i].Id]
			batchReplyUpdateRWMutex.RUnlock()
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
	}

	updatedMessages, errUM := json.Marshal(existingMessages)
	if errUM != nil {
		util.Loggify(errUM)
		return
	}
	success, errU := Update(&batchIdStruct.UserId, batchId, &updatedMessages, directory)
	if errU != nil || !success {
		util.Loggify(errU)
		return
	}
}

func SaveDirectMessageUpdates() {
	batchUpdateRWMutex.RLock()
	batchIds := make([]string, 0, len(newUpdates))
	for k := range newUpdates {
		batchIds = append(batchIds, k)
	}
	batchUpdateRWMutex.RUnlock()
	for _, batchId := range batchIds {
		//Acure a user specific mutex lock befor sending the message
		batchUpdateMutex.Lock()
		batchMutex := batchUpdateMutexis[batchId]
		if batchMutex == nil {
			batchUpdateMutexis[batchId] = &sync.Mutex{}
			batchMutex = batchUpdateMutexis[batchId]
		}
		batchUpdateMutex.Unlock()
		batchMutex.Lock()
		saveMessageUpdatesInternal(&batchId)
		batchUpdateRWMutex.Lock()
		delete(newUpdates, batchId)
		batchUpdateRWMutex.Unlock()
		batchMutex.Unlock()
	}
}

func saveMessageUpdatesInternal(batchId *string) {
	batchIdStruct, errBI := structs.ParseId(batchId)
	if errBI != nil {
		util.Loggify(errBI)
		return
	}
	optionalValueSplit := strings.Split(*batchIdStruct.OptionalValue, "-")
	if len(optionalValueSplit) != 2 {
		util.Loggify(fmt.Errorf("not able to save message updates due to invalid optional value is set on id"))
		return
	}
	var directory *string
	if optionalValueSplit[0] == "direct" {
		directory = util.PointerString(filepath.Join(CHAT_PATH, optionalValueSplit[1]))
	} else {
		util.Loggify(fmt.Errorf("not able to save message updates due to invalid optional value message type is set on id"))
		return
	}
	reader, _, errGF := GetFile(&batchIdStruct.UserId, directory, batchId, &structs.Version{
		Major: 0,
		Minor: 0,
		Patch: 1,
	}, util.PointerBoolean(false))
	if errGF != nil {
		util.Loggify(errGF)
		return
	}
	existingMessagesData, errPD := io.ReadAll(*reader)
	if errPD != nil {
		util.Loggify(errPD)
		return
	}
	var existingMessages []structs.Message
	errUM := json.Unmarshal(existingMessagesData, &existingMessages)
	if errUM != nil {
		util.Loggify(errUM)
		return
	}
	batchUpdateRWMutex.RLock()
	batchMessageIds := make([]string, 0, len(newUpdates[*batchId]))
	for k := range newUpdates[*batchId] {
		batchMessageIds = append(batchMessageIds, k)
	}
	batchUpdateRWMutex.RUnlock()
	for i := 0; i < len(existingMessages); i++ {
		if arrays.Contains(batchMessageIds, *existingMessages[i].Id) {
			batchUpdateRWMutex.RLock()
			messageUpdates := newUpdates[*batchId][*existingMessages[i].Id]
			batchUpdateRWMutex.RUnlock()
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
	}

	updatedMessages, errUM := json.Marshal(existingMessages)
	if errUM != nil {
		util.Loggify(errUM)
		return
	}
	success, errU := Update(&batchIdStruct.UserId, batchId, &updatedMessages, directory)
	if errU != nil || !success {
		util.Loggify(errU)
		return
	}
}

var CHAT_FILE_CHUNK_SIZE = 1024 * 1024
var CHAT_PATH = filepath.Join("__system__", "Chat")

func SaveDirectChatHistory() {
	newDirectMessageRWMutex.RLock()
	userIdKeys := make([]int64, 0, len(newDirectMessages))
	for k := range newDirectMessages {
		userIdKeys = append(userIdKeys, k)
	}
	newDirectMessageRWMutex.RUnlock()
	for _, userId := range userIdKeys {
		newDirectMessageMutex.Lock()
		userMutex := newDirectMessageMutexes[userId]
		if userMutex == nil {
			newDirectMessageMutexes[userId] = &sync.Mutex{}
			userMutex = newDirectMessageMutexes[userId]
		}
		newDirectMessageMutex.Unlock()
		userMutex.Lock()
		chatDirectory, errCh := GetChatDirectory(&userId)
		if errCh != nil {
			util.Loggify(errCh)
			userMutex.Unlock()
			continue
		}
		//Since the below code is reading a map, the reading opperation has to be synchronized
		//because in golang map read writes have to be synchronized
		newDirectMessageRWMutex.RLock()
		userToUserMessages := newDirectMessages[userId]
		userToUserMessagesKeys := make([]string, 0, len(userToUserMessages))
		for k := range userToUserMessages {
			userToUserMessagesKeys = append(userToUserMessagesKeys, k)
		}
		newDirectMessageRWMutex.RUnlock()

		for _, userIdPairString := range userToUserMessagesKeys {
			newDirectMessageRWMutex.RLock()
			messages := userToUserMessages[userIdPairString]
			newDirectMessageRWMutex.RUnlock()

			userIdPair := strings.Split(userIdPairString, "_")
			var to *int64
			if userIdPair[0] == strconv.FormatInt(userId, 10) {
				toTemp, _ := strconv.ParseInt(userIdPair[1], 10, 64)
				to = &toTemp
			} else {
				toTemp, _ := strconv.ParseInt(userIdPair[0], 10, 64)
				to = &toTemp
			}
			optionalValue := util.PointerString("direct-" + strconv.FormatInt(*to, 10))

			historyFile := filepath.Join(*chatDirectory, "directs", strconv.FormatInt(*to, 10), "history")
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
				if len(history.HistoryItems) == 0 || history.HistoryItems[len(history.HistoryItems)-1].Full {
					messagesBytes, errHB := json.Marshal(messages)
					if errHB != nil {
						util.Loggify(errHB)
						continue
					}
					size := int64(len(messagesBytes))
					date := datetime.FormatDate(time.Now())
					success, fileMetadata, errU := Upload(&userId, nil, &messagesBytes, util.PointerString(date), util.PointerString("application/json"), &size, util.PointerString(filepath.Join(CHAT_PATH, strconv.FormatInt(*to, 10))), optionalValue)
					if success && errU == nil {
						history.HistoryItems = append(history.HistoryItems, structs.ChatHistoryItem{
							Id:   fileMetadata.Id,
							Full: false,
						})

						historyBytes, errHB := json.Marshal(history)
						if errHB != nil {
							util.Loggify(errHB)
							continue
						}
						errR := os.Remove(historyFile)
						if errR != nil {
							util.Loggify(errR)
							continue
						}
						errH := os.WriteFile(historyFile, historyBytes, 0644)
						if errH != nil {
							util.Loggify(errH)
							continue
						}
						//Send update to all active participants for the new batch added
						//All active ui users should add the batch id to the chat item in client side
						var userUpdate interface{} = structs.UserMessageUpdate{
							MessageType:    util.PointerString("thread"),
							UpdateType:     util.PointerString("newBatch"),
							OpperationType: util.PointerString("newBatch"),
							From:           util.PointerString(strconv.FormatInt(userId, 10)),
							To:             util.PointerString(strconv.FormatInt(*to, 10)),
							Value:          fileMetadata.Id,
						}
						PutUpdate(&userId, util.PointerString("chat"), &userUpdate)
					} else if errU != nil {
						util.Loggify(errU)
						continue
					}
				} else {
					id := history.HistoryItems[len(history.HistoryItems)-1].Id
					reader, _, errGF := GetFile(&userId, util.PointerString(filepath.Join(CHAT_PATH, strconv.FormatInt(*to, 10))), id, &structs.Version{
						Major: 0,
						Minor: 0,
						Patch: 1,
					}, util.PointerBoolean(false))
					if errGF != nil {
						util.Loggify(errGF)
						continue
					}
					existingMessagesData, errPD := io.ReadAll(*reader)
					if errPD != nil {
						util.Loggify(errPD)
						continue
					}
					var existingMessages []structs.Message
					errUM := json.Unmarshal(existingMessagesData, &existingMessages)
					if errUM != nil {
						util.Loggify(errUM)
						continue
					}
					mergedMessages := append(existingMessages, messages...)
					updatedMessages, errUM := json.Marshal(mergedMessages)
					if errUM != nil {
						util.Loggify(errUM)
						continue
					}
					success, errU := Update(&userId, id, &updatedMessages, util.PointerString(filepath.Join(CHAT_PATH, strconv.FormatInt(*to, 10))))
					if errU != nil || !success {
						util.Loggify(errU)
						continue
					}
					if len(updatedMessages) >= CHAT_FILE_CHUNK_SIZE {
						history.HistoryItems[len(history.HistoryItems)-1].Full = true
						historyBytes, errHB := json.Marshal(history)
						if errHB != nil {
							util.Loggify(errHB)
							continue
						}
						errR := os.Remove(historyFile)
						if errR != nil {
							util.Loggify(errR)
							continue
						}
						errH := os.WriteFile(historyFile, historyBytes, 0644)
						if errH != nil {
							util.Loggify(errH)
							continue
						}
					}
				}
			}
			newDirectMessageRWMutex.Lock()
			delete(newDirectMessages[userId], userIdPairString)
			newDirectMessageRWMutex.Unlock()
		}
		userMutex.Unlock()
	}
}

var newRepliesCS = csmap.Create[string, []structs.Message](
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

func SendReplyMessage(userID *int64, sessionId *string, BatchId *string, From *string, To *string, RepliesToMessageId *string, RepliesToMessageBatchId *string, RepliesToMessageBatchId1 *string, RepliesToMessageType *string, Id *string, Text *string, Attachments *[]structs.ChatAttachment, Participants *[]string) (success bool, message *structs.Message, err error) {
	user := business.GetUser(userID, true)
	if user == nil || *user.AmperId == 0 || *user.AmperId != *business.AmperId() {
		return false, nil, fmt.Errorf("user is not allocated to this amper instance %d", *business.AmperId())
	}
	switch *RepliesToMessageType {
	case "direct":
		return sendDirecReplyMessage(userID, sessionId, BatchId, From, To, RepliesToMessageId, RepliesToMessageBatchId, RepliesToMessageBatchId1, RepliesToMessageType, Id, Text, Attachments, Participants)
	case "thread":
		return sendThreadReplyMessage(userID, sessionId, BatchId, From, To, RepliesToMessageId, RepliesToMessageBatchId, RepliesToMessageBatchId1, RepliesToMessageType, Id, Text, Attachments, Participants)
	case "channel":
		return sendChannelReplyMessage(userID, sessionId, BatchId, From, To, RepliesToMessageId, RepliesToMessageBatchId, RepliesToMessageBatchId1, RepliesToMessageType, Id, Text, Attachments, Participants)
	}
	return false, nil, nil
}

func sendDirecReplyMessage(userID *int64, sessionId *string, BatchId *string, From *string, To *string, RepliesToMessageId *string, RepliesToMessageBatchId *string, RepliesToMessageBatchId1 *string, RepliesToMessageType *string, Id *string, Text *string, Attachments *[]structs.ChatAttachment, Participants *[]string) (success bool, message *structs.Message, err error) {
	if BatchId == nil {
		messagesBytes := []byte("[]")
		size := int64(len(messagesBytes))
		date := datetime.FormatDate(time.Now())
		optionalValue := util.PointerString("direct-" + *To)
		success, fileMetadata, errU := Upload(userID, nil, &messagesBytes, util.PointerString(date), util.PointerString("application/json"), &size, util.PointerString(filepath.Join(CHAT_PATH, *To)), optionalValue)
		if success && errU == nil {
			BatchId = fileMetadata.Id
		} else {
			util.Loggify(errU)
			return false, nil, fmt.Errorf("not able to reserve a batch file, try again later")
		}
		success, errUC := UpdateChat(userID, sessionId, RepliesToMessageBatchId, RepliesToMessageBatchId1, RepliesToMessageId, RepliesToMessageType, util.PointerString("replyInitialisation"), util.PointerString("replyInitialisation"), BatchId, From, To)
		if errUC != nil || !success {
			util.Loggify(errUC)
			return false, nil, fmt.Errorf("not able to update the reserved batch file, try again later")
		}
	}
	success, errUC := UpdateChat(userID, sessionId, RepliesToMessageBatchId, RepliesToMessageBatchId1, RepliesToMessageId, RepliesToMessageType, util.PointerString("reply"), util.PointerString("reply"), From, From, To)
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

	newRepliesCS.StoreCompute(*BatchId, func(value []structs.Message) []structs.Message {
		if value != nil {
			return append(value, *message)
		} else {
			return []structs.Message{*message}
		}
	})

	var messageUpdate = structs.UserMessageUpdate{
		Message:        message,
		MessageType:    util.PointerString("direct"),
		UpdateType:     util.PointerString("newMessage"),
		OpperationType: util.PointerString("newMessage"),
		From:           From,
		To:             To,
	}
	SendUpdateToAllParticipants(userID, sessionId, Participants, &messageUpdate)
	return true, message, nil
}

func SendMessage(userID *int64, sessionId *string, from *string, to *string, participants *[]int64, Type *string, id *string, text *string, Attachments *[]structs.ChatAttachment) (success bool, message *structs.Message, err error) {
	user := business.GetUser(userID, true)
	if user == nil || *user.AmperId == 0 || *user.AmperId != *business.AmperId() {
		return false, nil, fmt.Errorf("user is not allocated to this amper instance %d", *business.AmperId())
	}

	switch *Type {
	case "direct":
		return sendDirectMessage(userID, sessionId, from, to, Type, id, text, Attachments)
	case "thread":
		return sendThreadMessage(userID, sessionId, from, to, participants, Type, id, text, Attachments)
	case "channel":
		return sendChannelMessage(userID, sessionId, from, to, Type, id, text, Attachments)
	}

	return false, nil, nil
}

func sendDirectMessage(userID *int64, sessionId *string, from *string, to *string, Type *string, id *string, text *string, Attachments *[]structs.ChatAttachment) (success bool, message *structs.Message, err error) {
	fromInt64, errFI := strconv.ParseInt(*from, 10, 64)
	if errFI != nil {
		return false, nil, fmt.Errorf("the from parameter supplied doesn't resolve to a valid identifying number")
	}
	toInt64, errTI := strconv.ParseInt(*to, 10, 64)
	if errTI != nil {
		return false, nil, fmt.Errorf("the to parameter supplied doesn't resolve to a valid identifying number")
	}
	if *userID != fromInt64 {
		return false, nil, fmt.Errorf("sending user must be the one currently logged in to the system")
	}

	//Acure a user specific mutex lock befor sending the message
	newDirectMessageMutex.Lock()
	userMutex := newDirectMessageMutexes[fromInt64]
	if userMutex == nil {
		newDirectMessageMutexes[fromInt64] = &sync.Mutex{}
		userMutex = newDirectMessageMutexes[fromInt64]
	}
	newDirectMessageMutex.Unlock()
	userMutex.Lock()
	defer userMutex.Unlock()

	optionalValue := util.PointerString("direct-" + *to)
	//Initilize user direct chat configuration if not exist
	chatDirectory, errCh := GetChatDirectory(&fromInt64)
	if errCh != nil {
		return false, nil, fmt.Errorf("not able to locate chat directory")
	}
	var batchFileId *string
	historyFile := filepath.Join(*chatDirectory, "directs", strconv.FormatInt(toInt64, 10), "history")
	if _, err := os.Stat(historyFile); errors.Is(err, os.ErrNotExist) {
		messagesBytes := []byte("[]")
		size := int64(len(messagesBytes))
		date := datetime.FormatDate(time.Now())
		success, fileMetadata, errU := Upload(userID, nil, &messagesBytes, util.PointerString(date), util.PointerString("application/json"), &size, util.PointerString(filepath.Join(CHAT_PATH, *to)), optionalValue)
		if success && errU == nil {
			batchFileId = fileMetadata.Id
			history := structs.ChatHistory{
				From: from,
				To:   to,
				HistoryItems: []structs.ChatHistoryItem{{
					Id:   fileMetadata.Id,
					Full: false,
				}},
				LastUpdateTime: time.Now().UnixMilli(),
			}
			historyBytes, errHB := json.Marshal(history)
			if errHB != nil {
				util.Loggify(errHB)
				return false, nil, fmt.Errorf("not able to initilize the history file due to json error")
			}

			errChD := os.MkdirAll(filepath.Join(*chatDirectory, "directs", strconv.FormatInt(toInt64, 10)), os.ModePerm)
			if errChD != nil {
				util.Loggify(errChD)
				return false, nil, fmt.Errorf("not able to initilize the history file directory")
			}

			errH := os.WriteFile(historyFile, historyBytes, 0644)
			if errH != nil {
				util.Loggify(errH)
				return false, nil, fmt.Errorf("not able to initilize the history file")
			}
		} else {
			util.Loggify(errU)
			return false, nil, fmt.Errorf("not able to reserve a history item")
		}
	} else {
		historyFileBytes, errHFB := os.ReadFile(historyFile)
		if errHFB != nil {
			util.Loggify(errHFB)
			return false, nil, fmt.Errorf("not able to read the history file")
		}
		history := &structs.ChatHistory{}
		errH := json.Unmarshal(historyFileBytes, history)
		if errH != nil {
			util.Loggify(errH)
			return false, nil, fmt.Errorf("not able to parse the history file")
		}
		if len(history.HistoryItems) == 0 || history.HistoryItems[len(history.HistoryItems)-1].Full {
			messagesBytes := []byte("[]")
			size := int64(len(messagesBytes))
			date := datetime.FormatDate(time.Now())
			success, fileMetadata, errU := Upload(userID, nil, &messagesBytes, util.PointerString(date), util.PointerString("application/json"), &size, util.PointerString(filepath.Join(CHAT_PATH, *to)), optionalValue)
			if success && errU == nil {
				batchFileId = fileMetadata.Id
				history.HistoryItems = append(history.HistoryItems, structs.ChatHistoryItem{
					Id:   fileMetadata.Id,
					Full: false,
				})
				history.LastUpdateTime = time.Now().UnixMilli()
				historyBytes, errHB := json.Marshal(history)
				if errHB != nil {
					util.Loggify(errHB)
					return false, nil, fmt.Errorf("not able to initilize the history file due to json error 1")
				}

				errR := os.Remove(historyFile)
				if errR != nil {
					util.Loggify(errR)
					return false, nil, fmt.Errorf("not able to remove the history file for rewriting")
				}
				errChD := os.MkdirAll(filepath.Join(*chatDirectory, "directs", strconv.FormatInt(toInt64, 10)), os.ModePerm)
				if errChD != nil {
					util.Loggify(errChD)
					return false, nil, fmt.Errorf("not able to initilize the history file directory 1")
				}

				errH := os.WriteFile(historyFile, historyBytes, 0644)
				if errH != nil {
					util.Loggify(errH)
					return false, nil, fmt.Errorf("not able to initilize the history file 1")
				}
			} else {
				util.Loggify(errU)
				return false, nil, fmt.Errorf("not able to reserve a history item 1")
			}
		} else {
			batchFileId = history.HistoryItems[len(history.HistoryItems)-1].Id

			history.LastUpdateTime = time.Now().UnixMilli()
			historyBytes, errHB := json.Marshal(history)
			if errHB != nil {
				util.Loggify(errHB)
				return false, nil, fmt.Errorf("not able to initilize the history file due to json error 1")
			}

			errR := os.Remove(historyFile)
			if errR != nil {
				util.Loggify(errR)
				return false, nil, fmt.Errorf("not able to remove the history file for rewriting")
			}
			errChD := os.MkdirAll(filepath.Join(*chatDirectory, "directs", strconv.FormatInt(toInt64, 10)), os.ModePerm)
			if errChD != nil {
				util.Loggify(errChD)
				return false, nil, fmt.Errorf("not able to initilize the history file directory 1")
			}

			errH := os.WriteFile(historyFile, historyBytes, 0644)
			if errH != nil {
				util.Loggify(errH)
				return false, nil, fmt.Errorf("not able to initilize the history file 1")
			}
		}
	}

	if batchFileId == nil {
		return false, nil, fmt.Errorf("not able to initilize the batch for the message")
	}

	toUser := business.GetUser(&toInt64, true)

	instance := business.GetAmperInstance(*toUser.AmperId)
	messageTime := time.Now().UnixMilli()
	success, value, errPS := DedicatedCallWithRetry(userID, sessionId, map[string]string{
		"amperInstance": "chat/receive",
	}, map[string]interface{}{
		"from":        from,
		"to":          to,
		"type":        Type,
		"id":          id,
		"text":        text,
		"time":        messageTime,
		"batchId":     batchFileId,
		"attachments": Attachments,
	}, instance)
	if !success || errPS != nil {
		util.Loggify(errPS)
		return false, nil, fmt.Errorf("not able to send the message to the remote user server")
	}

	message = &structs.Message{
		From:        from,
		To:          to,
		Id:          id,
		Text:        text,
		DateTime:    messageTime,
		BatchId:     batchFileId,
		BatchId1:    value,
		Attachments: Attachments,
	}
	newDirectMessageRWMutex.Lock()
	newMessageKey := getMessageKey(fromInt64, toInt64)
	if newDirectMessages[fromInt64] == nil {
		newDirectMessages[fromInt64] = make(map[string][]structs.Message)
	}
	newDirectMessages[fromInt64][newMessageKey] = append(newDirectMessages[fromInt64][newMessageKey], *message)
	newDirectMessageRWMutex.Unlock()
	return true, message, nil
}

func getMessageKey(from int64, to int64) string {
	if from > to {
		return fmt.Sprintf("%d_%d", from, to)
	} else {
		return fmt.Sprintf("%d_%d", to, from)
	}
}
func ReceiveMessage(userID *int64, from *string, to *string, Type *string, id *string, text *string, time int64, batchId *string, Attachments *[]structs.ChatAttachment) (success bool, batchIdResult *string, err error) {
	switch *Type {
	case "direct":
		return receiveDirectMessage(userID, from, to, id, text, time, batchId, Attachments)
	}

	return false, nil, nil
}

func receiveDirectMessage(userID *int64, from *string, to *string, id *string, text *string, messageTime int64, batchId *string, Attachments *[]structs.ChatAttachment) (success bool, batchIdResult *string, err error) {
	fromInt64, errFI := strconv.ParseInt(*from, 10, 64)
	if errFI != nil {
		return false, nil, fmt.Errorf("the from parameter supplied doesn't resolve to a valid identifying number")
	}
	toInt64, errTI := strconv.ParseInt(*to, 10, 64)
	if errTI != nil {
		return false, nil, fmt.Errorf("the to parameter supplied doesn't resolve to a valid identifying number")
	}
	user := business.GetUser(&toInt64, true)
	if user == nil || *user.AmperId == 0 || *user.AmperId != *business.AmperId() {
		return false, nil, fmt.Errorf("not able to receive, user is not allocated to this amper instance %d", *business.AmperId())
	}
	if *userID != fromInt64 {
		return false, nil, fmt.Errorf("sending user must be the one currently logged in to the system")
	}
	newDirectMessageMutex.Lock()
	//Acure a user specific mutex lock befor sending the message
	userMutex := newDirectMessageMutexes[toInt64]
	if userMutex == nil {
		//Acuire a global lock for initilizing the user mutex
		newDirectMessageMutexes[toInt64] = &sync.Mutex{}
		userMutex = newDirectMessageMutexes[toInt64]
	}
	newDirectMessageMutex.Unlock()
	userMutex.Lock()
	defer userMutex.Unlock()

	//Initilize user direct chat configuration if not exist
	chatDirectory, errCh := GetChatDirectory(&toInt64)
	if errCh != nil {
		return false, nil, fmt.Errorf("not able to locate chat directory 1")
	}
	optionalValue := util.PointerString("direct-" + *from)
	var batchFileId *string
	historyFile := filepath.Join(*chatDirectory, "directs", strconv.FormatInt(fromInt64, 10), "history")
	var chatHistory structs.ChatHistory
	if _, err := os.Stat(historyFile); errors.Is(err, os.ErrNotExist) {
		messagesBytes := []byte("[]")
		size := int64(len(messagesBytes))
		date := datetime.FormatDate(time.Now())
		success, fileMetadata, errU := Upload(&toInt64, nil, &messagesBytes, util.PointerString(date), util.PointerString("application/json"), &size, util.PointerString(filepath.Join(CHAT_PATH, *from)), optionalValue)
		if success && errU == nil {
			batchFileId = fileMetadata.Id
			chatHistory = structs.ChatHistory{
				From: to,
				To:   from,
				HistoryItems: []structs.ChatHistoryItem{{
					Id:   fileMetadata.Id,
					Full: false,
				}},
				LastUpdateTime: time.Now().UnixMilli(),
				UnreadMessages: 1,
			}
			historyBytes, errHB := json.Marshal(chatHistory)
			if errHB != nil {
				util.Loggify(errHB)
				return false, nil, fmt.Errorf("not able to initilize the history file due to json error")
			}

			errChD := os.MkdirAll(filepath.Join(*chatDirectory, "directs", strconv.FormatInt(fromInt64, 10)), os.ModePerm)
			if errChD != nil {
				util.Loggify(errChD)
				return false, nil, fmt.Errorf("not able to initilize the history file directory")
			}

			errH := os.WriteFile(historyFile, historyBytes, 0644)
			if errH != nil {
				util.Loggify(errH)
				return false, nil, fmt.Errorf("not able to initilize the history file")
			}
		} else {
			util.Loggify(errU)
			return false, nil, fmt.Errorf("not able to reserve a history item 1")
		}
	} else {
		historyFileBytes, errHFB := os.ReadFile(historyFile)
		if errHFB != nil {
			util.Loggify(errHFB)
			return false, nil, fmt.Errorf("not able to read the history file")
		}
		errH := json.Unmarshal(historyFileBytes, &chatHistory)
		if errH != nil {
			util.Loggify(errH)
			return false, nil, fmt.Errorf("not able to parse the history file")
		}
		if len(chatHistory.HistoryItems) == 0 || chatHistory.HistoryItems[len(chatHistory.HistoryItems)-1].Full {
			messagesBytes := []byte("[]")
			size := int64(len(messagesBytes))
			date := datetime.FormatDate(time.Now())
			success, fileMetadata, errU := Upload(&toInt64, nil, &messagesBytes, util.PointerString(date), util.PointerString("application/json"), &size, util.PointerString(filepath.Join(CHAT_PATH, *from)), optionalValue)
			if success && errU == nil {
				batchFileId = fileMetadata.Id
				chatHistory.HistoryItems = append(chatHistory.HistoryItems, structs.ChatHistoryItem{
					Id:   fileMetadata.Id,
					Full: false,
				})
				chatHistory.LastUpdateTime = time.Now().UnixMilli()
				chatHistory.UnreadMessages++
				historyBytes, errHB := json.Marshal(chatHistory)
				if errHB != nil {
					util.Loggify(errHB)
					return false, nil, fmt.Errorf("not able to initilize the history file due to json error 2")
				}

				errR := os.Remove(historyFile)
				if errR != nil {
					util.Loggify(errR)
					return false, nil, fmt.Errorf("not able to remove the history file for rewriting 1")
				}
				errChD := os.MkdirAll(filepath.Join(*chatDirectory, "directs", strconv.FormatInt(toInt64, 10)), os.ModePerm)
				if errChD != nil {
					util.Loggify(errChD)
					return false, nil, fmt.Errorf("not able to initilize the history file directory 2")
				}

				errH := os.WriteFile(historyFile, historyBytes, 0644)
				if errH != nil {
					util.Loggify(errH)
					return false, nil, fmt.Errorf("not able to initilize the history file 2")
				}
			} else {
				util.Loggify(errU)
				return false, nil, fmt.Errorf("not able to reserve a history item 2")
			}
		} else {
			batchFileId = chatHistory.HistoryItems[len(chatHistory.HistoryItems)-1].Id
			chatHistory.LastUpdateTime = time.Now().UnixMilli()
			chatHistory.UnreadMessages++
			historyBytes, errHB := json.Marshal(chatHistory)
			if errHB != nil {
				util.Loggify(errHB)
				return false, nil, fmt.Errorf("not able to initilize the history file due to json error 2")
			}

			errR := os.Remove(historyFile)
			if errR != nil {
				util.Loggify(errR)
				return false, nil, fmt.Errorf("not able to remove the history file for rewriting 1")
			}
			errChD := os.MkdirAll(filepath.Join(*chatDirectory, "directs", strconv.FormatInt(toInt64, 10)), os.ModePerm)
			if errChD != nil {
				util.Loggify(errChD)
				return false, nil, fmt.Errorf("not able to initilize the history file directory 2")
			}

			errH := os.WriteFile(historyFile, historyBytes, 0644)
			if errH != nil {
				util.Loggify(errH)
				return false, nil, fmt.Errorf("not able to initilize the history file 2")
			}
		}
	}
	message := structs.Message{
		From:        from,
		To:          to,
		Id:          id,
		Text:        text,
		DateTime:    messageTime,
		BatchId:     batchFileId,
		BatchId1:    batchId,
		Attachments: Attachments,
	}
	newDirectMessageRWMutex.Lock()
	newMessageKey := getMessageKey(fromInt64, toInt64)
	if newDirectMessages[toInt64] == nil {
		newDirectMessages[toInt64] = make(map[string][]structs.Message)
	}
	newDirectMessages[toInt64][newMessageKey] = append(newDirectMessages[toInt64][newMessageKey], message)
	newDirectMessageRWMutex.Unlock()

	//Inform the user for receiving an update
	var users []structs.User = []structs.User{*business.GetUser(&fromInt64, true)}
	var userUpdate interface{} = structs.UserMessageUpdate{
		Message:        &message,
		MessageType:    util.PointerString("direct"),
		UpdateType:     util.PointerString("newMessage"),
		OpperationType: util.PointerString("newMessage"),
		From:           from,
		To:             to,
		ChatHistory:    &chatHistory,
		Users:          &users,
	}
	PutUpdate(&toInt64, util.PointerString("chat"), &userUpdate)
	return true, batchFileId, nil
}

func GetChatDirectory(userId *int64) (result *string, err error) {
	instance := business.GetAmperInstance(*business.AmperId())
	rootDriveDirectory := instance.Directory
	if rootDriveDirectory == nil || len(*rootDriveDirectory) < 1 {
		return result, fmt.Errorf("empty directory found, amper instance %d is not configured for a directory", *business.AmperId())
	}
	result = util.PointerString(filepath.Join(*rootDriveDirectory, strconv.FormatInt(*userId, 10), "chat"))
	errUD := os.MkdirAll(*result, os.ModePerm)
	if errUD != nil && !errors.Is(err, os.ErrExist) {
		util.Loggify(errUD)
		return result, fmt.Errorf("not able to locate the user's active directory in chat '%s', please contect the support", *result)
	}
	return result, nil
}

func GetChatState(userID *int64, sessionId *string) (result *structs.ChatState, err error) {
	result = &structs.ChatState{}
	chatDirectory, errCh := GetChatDirectory(userID)
	if errCh != nil {
		return nil, fmt.Errorf("not able to locate chat directory")
	}

	//Read direct chat items
	directsDirectory := filepath.Join(*chatDirectory, "directs")
	entries, errE := os.ReadDir(directsDirectory)
	if errE != nil {
		util.Loggify(errE)
	}
	var directs []structs.ChatDirectItem = make([]structs.ChatDirectItem, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			historyFile := filepath.Join(directsDirectory, entry.Name(), "history")
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
			toInt64, _ := strconv.ParseInt(*history.To, 10, 64)
			user := business.GetUser(&toInt64, true)

			chatDirectItem := structs.ChatDirectItem{
				User:        user,
				ChatHistory: history,
			}
			directs = append(directs, chatDirectItem)
		}
	}
	result.Directs = &directs

	//Read thread chat items
	threadsDirectory := filepath.Join(*chatDirectory, "threads")
	threadsEntries, errE := os.ReadDir(threadsDirectory)
	if errE != nil {
		util.Loggify(errE)
	}
	var threads []structs.ChatThreadItem = make([]structs.ChatThreadItem, 0)
	var instanceThreadIds map[int64][]string = map[int64][]string{}
	var threadToHistoryMap map[string]structs.ChatThreadHistory = make(map[string]structs.ChatThreadHistory)
	for _, threadEntry := range threadsEntries {
		historyFile := filepath.Join(threadsDirectory, threadEntry.Name(), "history")
		historyFileBytes, errHFB := os.ReadFile(historyFile)
		if errHFB != nil {
			util.Loggify(errHFB)
			continue
		}
		history := &structs.ChatThreadHistory{}
		errH := json.Unmarshal(historyFileBytes, history)
		if errH != nil {
			util.Loggify(errH)
			continue
		}

		threadIdSplit := strings.Split(*history.ThreadId, "_")
		if len(threadIdSplit) == 2 {
			threadToHistoryMap[*history.ThreadId] = *history
			threadInstanceId, errTII := strconv.ParseInt(threadIdSplit[0], 10, 64)
			if errTII == nil {
				if instanceThreadIds[threadInstanceId] != nil {
					instanceThreadIds[threadInstanceId] = append(instanceThreadIds[threadInstanceId], *history.ThreadId)
				} else {
					instanceThreadIds[threadInstanceId] = []string{*history.ThreadId}
				}
			}
		}
	}
	for instanceId, threadIds := range instanceThreadIds {
		instance := business.GetAmperInstance(instanceId)
		successInstance, value, errPS := DedicatedCallWithRetry(userID, sessionId, map[string]string{
			"amperInstance": "chat/getThreads",
		}, map[string]interface{}{
			"threadIds": threadIds,
		}, instance)
		if errPS == nil && successInstance {
			chatHistories := []structs.ChatHistory{}
			json.Unmarshal([]byte(*value), &chatHistories)
			for i := 0; i < len(chatHistories); i++ {
				chatHistory := chatHistories[i]
				participants := []structs.User{}
				if chatHistory.Participants != nil {
					for _, participant := range *chatHistory.Participants {
						participants = append(participants, *business.GetUser(&participant, true))
					}
				}
				history := threadToHistoryMap[*chatHistory.To]
				chatHistory.UnreadMessages = history.UnreadMessages
				chatHistory.LastUpdateTime = history.LastUpdateTime
				chatThreadItem := structs.ChatThreadItem{
					Label:       history.Label,
					ChatHistory: &chatHistory,
					Users:       &participants,
				}
				threads = append(threads, chatThreadItem)
			}
		}
	}
	result.Threads = &threads

	//fetch chat channel state
	channelGroups, errCG := GetChatChannelState(userID, sessionId)
	if errCG == nil {
		result.ChannelChannelGroups = channelGroups
	} else {
		util.Loggify(errCG)
	}
	return result, nil
}

func FetchChatReplies(userID *int64, BatchId *string, MessageType *string) (result []structs.Message, participants map[int64]*structs.User, err error) {
	switch *MessageType {
	case "direct":
		return FetchChatRepliesDirect(userID, BatchId, MessageType)
	case "thread":
		return FetchChatRepliesThread(userID, BatchId, MessageType)
	case "channel":
		return FetchChatRepliesChannel(userID, BatchId, MessageType)
	}
	return []structs.Message{}, participants, nil
}

func FetchChatRepliesDirect(userID *int64, BatchId *string, MessageType *string) (result []structs.Message, participants map[int64]*structs.User, err error) {
	result = make([]structs.Message, 0)

	reader, errR := GetFileBody(userID, BatchId)
	if errR != nil {
		util.Loggify(errR)
		return nil, nil, fmt.Errorf("not able to retrieve the chat replies")
	}
	existingMessagesData, errPD := io.ReadAll(*reader)
	if errPD != nil {
		util.Loggify(errPD)
		return nil, nil, fmt.Errorf("not able to read the chat replies")
	}
	var existingMessages []structs.Message
	errUM := json.Unmarshal(existingMessagesData, &existingMessages)
	if errUM != nil {
		util.Loggify(errUM)
		return nil, nil, fmt.Errorf("not able to parse the chat replies")
	}

	var intermediateResult []structs.Message = make([]structs.Message, 0)
	for _, existingMessage := range existingMessages {
		if !existingMessage.Deleted {
			intermediateResult = append(intermediateResult, existingMessage)
		}
	}

	newRepliesToAppend, okNR := newRepliesCS.Load(*BatchId)
	if okNR {
		intermediateResult = append(intermediateResult, newRepliesToAppend...)
	}

	batchUpdateRWMutex.RLock()
FirstLoop:
	for i := 0; i < len(intermediateResult); i++ {
		if newReplyUpdates[*BatchId][*intermediateResult[i].Id] != nil {
			messageUpdates := newReplyUpdates[*BatchId][*intermediateResult[i].Id]
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
		}
		result = append(result, intermediateResult[i])
	}
	batchUpdateRWMutex.RUnlock()
	return result, participants, nil
}

func FetchChat(userID *int64, from *string, to *string, Type *string, id *string, includeLatest bool) (result []structs.Message, participants map[int64]*structs.User, err error) {
	switch *Type {
	case "direct":
		return fetchChatDirect(userID, from, to, id, includeLatest)
	case "thread":
		return fetchChatThread(userID, from, to, id, includeLatest)
	case "channel":
		return fetchChatChannel(userID, from, to, id, includeLatest)
	}

	return result, nil, nil
}

func fetchChatDirect(userID *int64, from *string, to *string, id *string, includeLatest bool) (result []structs.Message, participants map[int64]*structs.User, err error) {
	result = make([]structs.Message, 0)
	fromInt64, errFI := strconv.ParseInt(*from, 10, 64)
	if errFI != nil {
		return nil, nil, fmt.Errorf("not able to fetch, the from parameter supplied doesn't resolve to a valid identifying number")
	}
	toInt64, errTI := strconv.ParseInt(*to, 10, 64)
	if errTI != nil {
		return nil, nil, fmt.Errorf("not able to fetch, the to parameter supplied doesn't resolve to a valid identifying number")
	}
	if *userID != fromInt64 {
		return nil, nil, fmt.Errorf("not able to fetch, sending user must be the one currently logged in to the system")
	}

	//Acure a user specific mutex lock befor sending the message
	newDirectMessageMutex.Lock()
	userMutex := newDirectMessageMutexes[*userID]
	if userMutex == nil {
		//Acuire a global lock for initilizing the user mutex
		newDirectMessageMutexes[*userID] = &sync.Mutex{}
		userMutex = newDirectMessageMutexes[*userID]
	}
	newDirectMessageMutex.Unlock()
	userMutex.Lock()
	defer userMutex.Unlock()

	batchUpdateMutex.Lock()
	//Acure a batch specific mutex lock befor sending the message
	batchMutex := batchUpdateMutexis[*id]
	if batchMutex == nil {
		//Acuire a global lock for initilizing the user mutex
		batchUpdateMutexis[*id] = &sync.Mutex{}
		batchMutex = batchUpdateMutexis[*id]
	}
	batchUpdateMutex.Unlock()
	batchMutex.Lock()
	defer batchMutex.Unlock()

	var intermediateResult []structs.Message = make([]structs.Message, 0)
	reader, _, errGF := GetFile(userID, util.PointerString(filepath.Join(CHAT_PATH, strconv.FormatInt(toInt64, 10))), id, &structs.Version{
		Major: 0,
		Minor: 0,
		Patch: 1,
	}, util.PointerBoolean(false))
	if errGF == nil {
		existingMessagesData, errPD := io.ReadAll(*reader)
		if errPD == nil {
			var existingMessages []structs.Message
			errUM := json.Unmarshal(existingMessagesData, &existingMessages)
			if errUM == nil {
				for _, existingMessage := range existingMessages {
					if !existingMessage.Deleted {
						intermediateResult = append(intermediateResult, existingMessage)
					}
				}
			} else {
				util.Loggify(errUM)
			}
		} else {
			util.Loggify(errPD)
		}
	} else {
		util.Loggify(errGF)
	}

	if includeLatest {
		newDirectMessageRWMutex.RLock()
		if newDirectMessages[fromInt64] != nil {
			newMessageKey := getMessageKey(fromInt64, toInt64)
			newMessages := newDirectMessages[fromInt64][newMessageKey]
			intermediateResult = append(intermediateResult, newMessages...)
		}
		newDirectMessageRWMutex.RUnlock()
	}
	batchUpdateRWMutex.RLock()
FirstLoop:
	for i := 0; i < len(intermediateResult); i++ {
		if newUpdates[*id][*intermediateResult[i].Id] != nil {
			messageUpdates := newUpdates[*id][*intermediateResult[i].Id]
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
		}
		result = append(result, intermediateResult[i])
	}
	batchUpdateRWMutex.RUnlock()
	return result, nil, nil
}

func UpdateChatReply(userID *int64, sessionId *string, BatchId *string, MessageId *string, MessageType *string, UpdateType *string, OpperationType *string, Value *string, From *string, To *string, Participants *[]string) (success bool, err error) {
	switch *MessageType {
	case "direct":
		return UpdateChatReplyDirect(userID, sessionId, BatchId, MessageId, MessageType, UpdateType, OpperationType, Value, From, To, Participants)
	case "thread":
		return UpdateChatReplyThread(userID, sessionId, BatchId, MessageId, MessageType, UpdateType, OpperationType, Value, From, To, Participants)
	case "channel":
		return UpdateChatReplyChannel(userID, sessionId, BatchId, MessageId, MessageType, UpdateType, OpperationType, Value, From, To, Participants)
	}
	return false, fmt.Errorf("wrong message type suplied")
}

// The below cach is designed to hold new updates for reply chat that have not been persisted to the key value store
var newReplyUpdates map[string]map[string][]structs.MessageUpdate = make(map[string]map[string][]structs.MessageUpdate)
var batchReplyUpdateMutexis map[string]*sync.Mutex = make(map[string]*sync.Mutex)
var batchReplyUpdateMutex sync.Mutex
var batchReplyUpdateRWMutex sync.RWMutex

func UpdateChatReplyDirect(userID *int64, sessionId *string, BatchId *string, MessageId *string, MessageType *string, UpdateType *string, OpperationType *string, Value *string, From *string, To *string, Participants *[]string) (success bool, err error) {
	fromInt64, errFI := strconv.ParseInt(*From, 10, 64)
	if errFI != nil {
		return false, fmt.Errorf("the from parameter supplied doesn't resolve to a valid identifying number")
	}
	if *userID != fromInt64 {
		return false, fmt.Errorf("the reacting user must be the one logged in to the system")
	}
	batchReplyUpdateMutex.Lock()
	//Acure a user specific mutex lock befor sending the message
	batchReplyMutex := batchReplyUpdateMutexis[*BatchId]
	if batchReplyMutex == nil {
		//Acuire a global lock for initilizing the user mutex
		batchReplyUpdateMutexis[*BatchId] = &sync.Mutex{}
		batchReplyMutex = batchReplyUpdateMutexis[*BatchId]
	}
	batchReplyUpdateMutex.Unlock()
	batchReplyMutex.Lock()
	defer batchReplyMutex.Unlock()

	batchReplyUpdateRWMutex.Lock()
	if newReplyUpdates[*BatchId] == nil {
		newReplyUpdates[*BatchId] = make(map[string][]structs.MessageUpdate)
	}
	if newReplyUpdates[*BatchId][*MessageId] == nil {
		newReplyUpdates[*BatchId][*MessageId] = make([]structs.MessageUpdate, 0)
	}
	messageUpdate := structs.MessageUpdate{
		MessageType:    MessageType,
		UpdateType:     UpdateType,
		OpperationType: OpperationType,
		From:           From,
		To:             To,
		Value:          Value,
	}
	newReplyUpdates[*BatchId][*MessageId] = append(newReplyUpdates[*BatchId][*MessageId], messageUpdate)
	batchReplyUpdateRWMutex.Unlock()

	userMessageUpdate := structs.UserMessageUpdate{
		Message: &structs.Message{
			BatchId: BatchId,
			Id:      MessageId,
		},
		MessageType:    messageUpdate.MessageType,
		UpdateType:     messageUpdate.UpdateType,
		OpperationType: messageUpdate.OpperationType,
		From:           messageUpdate.From,
		To:             messageUpdate.To,
		Value:          messageUpdate.Value,
	}
	SendUpdateToAllParticipants(userID, sessionId, Participants, &userMessageUpdate)
	return true, nil
}

func SendUpdateToAllParticipants(userID *int64, sessionId *string, Participants *[]string, update *structs.UserMessageUpdate) {
	var instanceToParticipantsMap map[int64][]int64 = make(map[int64][]int64)
	for _, participant := range *Participants {
		participantInt64, errFI := strconv.ParseInt(participant, 10, 64)
		if errFI != nil {
			util.Loggify(errFI)
			continue
		}
		//no need to send update to itself
		if *userID == participantInt64 {
			continue
		}
		user := business.GetUser(&participantInt64, true)
		if instanceToParticipantsMap[*user.AmperId] == nil {
			instanceToParticipantsMap[*user.AmperId] = []int64{*user.ID}
		} else {
			instanceToParticipantsMap[*user.AmperId] = append(instanceToParticipantsMap[*user.AmperId], *user.ID)
		}
	}
	for instanceId, instanceParticipants := range instanceToParticipantsMap {
		instance := business.GetAmperInstance(instanceId)
		_, _, errPS := DedicatedCallWithRetry(userID, sessionId, map[string]string{
			"amperInstance": "updates/push",
		}, map[string]interface{}{
			"category":     "chatReply",
			"participants": instanceParticipants,
			"value":        update,
		}, instance)
		if errPS != nil {
			util.Loggify(errPS)
		}
	}

}
func UpdateChat(userID *int64, sessionId *string, BatchId *string, BatchId1 *string, MessageId *string, MessageType *string, UpdateType *string, OpperationType *string, Value *string, From *string, To *string) (success bool, err error) {
	switch *MessageType {
	case "direct":
		return updateChatDirect(userID, sessionId, BatchId, BatchId1, MessageId, MessageType, UpdateType, OpperationType, Value, From, To)
	case "thread":
		//This assumes the BatchId messages are stored in curent amper instance, so the batch id instance id is current instance
		return updateReceiveChatThread(userID, sessionId, BatchId, MessageId, MessageType, UpdateType, OpperationType, Value, From, To)
	case "channel":
		//This assumes the BatchId messages are stored in curent amper instance, so the batch id instance id is current instance
		return updateReceiveChatChannel(userID, sessionId, BatchId, MessageId, MessageType, UpdateType, OpperationType, Value, From, To)
	}
	return false, nil
}

func updateChatDirect(userID *int64, sessionId *string, BatchId *string, BatchId1 *string, MessageId *string, MessageType *string, UpdateType *string, OpperationType *string, Value *string, From *string, To *string) (success bool, err error) {
	fromInt64, errFI := strconv.ParseInt(*From, 10, 64)
	if errFI != nil {
		return false, fmt.Errorf("the from parameter supplied doesn't resolve to a valid identifying number")
	}
	if *userID != fromInt64 {
		return false, fmt.Errorf("the reacting user must be the one logged in to the system")
	}
	successInternal, errSI := updateChatInternal(userID, BatchId, MessageId, MessageType, UpdateType, OpperationType, Value, From, To)
	if !successInternal || errSI != nil {
		util.Loggify(errSI)
		return false, fmt.Errorf("not able to register the reaction in memory for internal batch")
	}
	//Update the second copy of the message corresponding to BatchId1
	if BatchId1 != nil {
		successRemote, errSR := updateChatRemote(userID, sessionId, BatchId1, MessageId, MessageType, UpdateType, OpperationType, Value, From, To)
		if !successRemote || errSR != nil {
			util.Loggify(errSR)
			return false, fmt.Errorf("not able to register the reaction in memory for remote batch")
		}
	}
	return true, nil
}

func updateChatRemote(userID *int64, sessionId *string, BatchId1 *string, MessageId *string, MessageType *string, UpdateType *string, OpperationType *string, Value *string, From *string, To *string) (success bool, err error) {
	batchId1Struct, errBI1 := structs.ParseId(BatchId1)
	if errBI1 != nil {
		return false, fmt.Errorf("not able to synch reaction with second batch due to invalid batchid provided")
	}
	toUser := business.GetUser(&batchId1Struct.UserId, true)

	instance := business.GetAmperInstance(*toUser.AmperId)

	success, _, errPS := DedicatedCallWithRetry(userID, sessionId, map[string]string{
		"amperInstance": "chat/updateReceive",
	}, map[string]interface{}{
		"batchId":        BatchId1,
		"messageId":      MessageId,
		"messageType":    MessageType,
		"updateType":     UpdateType,
		"opperationType": OpperationType,
		"value":          Value,
		"from":           From,
		"to":             To,
	}, instance)
	if !success || errPS != nil {
		util.Loggify(errPS)
		return false, fmt.Errorf("not able to send the message to the remote user server")
	}
	return true, nil
}

func updateChatInternal(userID *int64, BatchId *string, MessageId *string, MessageType *string, UpdateType *string, OpperationType *string, Value *string, From *string, To *string) (success bool, err error) {
	batchUpdateMutex.Lock()
	//Acure a user specific mutex lock befor sending the message
	batchMutex := batchUpdateMutexis[*BatchId]
	if batchMutex == nil {
		//Acuire a global lock for initilizing the user mutex
		batchUpdateMutexis[*BatchId] = &sync.Mutex{}
		batchMutex = batchUpdateMutexis[*BatchId]
	}
	batchUpdateMutex.Unlock()
	batchMutex.Lock()
	defer batchMutex.Unlock()

	batchUpdateRWMutex.Lock()
	if newUpdates[*BatchId] == nil {
		newUpdates[*BatchId] = make(map[string][]structs.MessageUpdate)
	}
	if newUpdates[*BatchId][*MessageId] == nil {
		newUpdates[*BatchId][*MessageId] = make([]structs.MessageUpdate, 0)
	}
	messageUpdate := structs.MessageUpdate{
		MessageType:    MessageType,
		UpdateType:     UpdateType,
		OpperationType: OpperationType,
		From:           From,
		To:             To,
		Value:          Value,
	}
	newUpdates[*BatchId][*MessageId] = append(newUpdates[*BatchId][*MessageId], messageUpdate)
	batchUpdateRWMutex.Unlock()
	return true, nil
}
func UpdateReceiveChat(userID *int64, sessionId *string, BatchId *string, MessageId *string, MessageType *string, UpdateType *string, OpperationType *string, Value *string, From *string, To *string) (success bool, err error) {
	switch *MessageType {
	case "direct":
		return UpdateReceiveChatDirect(userID, BatchId, MessageId, MessageType, UpdateType, OpperationType, Value, From, To)
	case "thread":
		return updateReceiveChatThread(userID, sessionId, BatchId, MessageId, MessageType, UpdateType, OpperationType, Value, From, To)
	}
	return false, nil
}

func UpdateReceiveChatDirect(userID *int64, BatchId *string, MessageId *string, MessageType *string, UpdateType *string, OpperationType *string, Value *string, From *string, To *string) (success bool, err error) {
	batchUpdateMutex.Lock()
	//Acure a user specific mutex lock befor sending the message
	batchMutex := batchUpdateMutexis[*BatchId]
	if batchMutex == nil {
		//Acuire a global lock for initilizing the user mutex
		batchUpdateMutexis[*BatchId] = &sync.Mutex{}
		batchMutex = batchUpdateMutexis[*BatchId]
	}
	batchUpdateMutex.Unlock()
	batchMutex.Lock()
	defer batchMutex.Unlock()

	batchUpdateRWMutex.Lock()
	if newUpdates[*BatchId] == nil {
		newUpdates[*BatchId] = make(map[string][]structs.MessageUpdate)
	}
	if newUpdates[*BatchId][*MessageId] == nil {
		newUpdates[*BatchId][*MessageId] = make([]structs.MessageUpdate, 0)
	}
	messageUpdate := structs.MessageUpdate{
		MessageType:    MessageType,
		UpdateType:     UpdateType,
		OpperationType: OpperationType,
		From:           From,
		To:             To,
		Value:          Value,
	}
	newUpdates[*BatchId][*MessageId] = append(newUpdates[*BatchId][*MessageId], messageUpdate)
	batchUpdateRWMutex.Unlock()
	//Inform the user for receiving an update
	var userUpdate interface{} = structs.UserMessageUpdate{
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
	toInt64, errTI := strconv.ParseInt(*To, 10, 64)
	if errTI != nil {
		util.Loggify(errTI)
	} else {
		PutUpdate(&toInt64, util.PointerString("chat"), &userUpdate)
	}
	return true, nil
}

func MarkChatUnread(userID *int64, to *string, chatType *string) (success bool, err error) {
	switch *chatType {
	case "direct":
		return markChatUnreadDirect(userID, to, chatType)
	case "thread":
		return markChatUnreadThread(userID, to, chatType)
	}
	return true, nil
}

func markChatUnreadDirect(userID *int64, to *string, chatType *string) (success bool, err error) {
	toInt64, errTI := strconv.ParseInt(*to, 10, 64)
	if errTI != nil {
		return false, fmt.Errorf("markChatUnreadDirect - the to parameter supplied doesn't resolve to a valid identifying number")
	}
	chatDirectory, errCh := GetChatDirectory(userID)
	if errCh != nil {
		return false, fmt.Errorf("markChatUnreadDirect - not able to locate chat directory")
	}
	historyFile := filepath.Join(*chatDirectory, "directs", strconv.FormatInt(toInt64, 10), "history")
	if _, err := os.Stat(historyFile); !errors.Is(err, os.ErrNotExist) {
		historyFileBytes, errHFB := os.ReadFile(historyFile)
		if errHFB != nil {
			util.Loggify(errHFB)
			return false, fmt.Errorf("markChatUnreadDirect - not able to read the history file")
		}
		history := &structs.ChatHistory{}
		errH := json.Unmarshal(historyFileBytes, history)
		if errH != nil {
			util.Loggify(errH)
			return false, fmt.Errorf("markChatUnreadDirect - not able to parse the history file")
		}
		history.UnreadMessages = 0
		historyBytes, errHB := json.Marshal(history)
		if errHB != nil {
			util.Loggify(errHB)
			return false, fmt.Errorf("markChatUnreadDirect - not able to initilize the history file due to json error 1")
		}

		errR := os.Remove(historyFile)
		if errR != nil {
			util.Loggify(errR)
			return false, fmt.Errorf("markChatUnreadDirect - not able to remove the history file for rewriting")
		}
		errChD := os.MkdirAll(filepath.Join(*chatDirectory, "directs", strconv.FormatInt(toInt64, 10)), os.ModePerm)
		if errChD != nil {
			util.Loggify(errChD)
			return false, fmt.Errorf("markChatUnreadDirect - not able to initilize the history file directory 1")
		}

		errH = os.WriteFile(historyFile, historyBytes, 0644)
		if errH != nil {
			util.Loggify(errH)
			return false, fmt.Errorf("markChatUnreadDirect - not able to initilize the history file 1")
		}
	}
	return true, nil
}
