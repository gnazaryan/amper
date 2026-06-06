package controller

import (
	"amper/auth/authorization"
	"amper/common/structs"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func ChatController(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result string) {
	pathSplit := strings.Split(r.URL.Path, "/")
	var resultStruct interface{}
	if len(pathSplit) > 2 {
		switch pathSplit[2] {
		case "send":
			resultStruct = sendMessage(userID, sessionId, w, r)
		case "sendReply":
			resultStruct = sendReplyMessage(userID, sessionId, w, r)
		case "receive":
			resultStruct = receiveMessage(userID, w, r)
		case "state":
			resultStruct = getChatState(userID, sessionId, w, r)
		case "fetch":
			resultStruct = fetchChat(userID, w, r)
		case "fetchReplies":
			resultStruct = fetchChatReplies(userID, w, r)
		case "update":
			resultStruct = updateChat(userID, sessionId, w, r)
		case "updateReceive":
			resultStruct = updateReceiveChat(userID, sessionId, w, r)
		case "updateReply":
			resultStruct = updateChatReply(userID, sessionId, w, r)
		case "markUnread":
			resultStruct = markChatUnread(userID, w, r)
		case "initThread":
			resultStruct = initThread(userID, w, r)
		case "receiveThread":
			resultStruct = receiveThread(userID, w, r)
		case "getThreads":
			resultStruct = getThreads(userID, w, r)
		case "updateThread":
			resultStruct = updateThread(userID, w, r)
		case "download":
			downloadAttachment(userID, w, r)
		case "createChatGroup":
			resultStruct = createChatGroup(userID, w, r)
		case "fetchChatChannelGroups":
			resultStruct = fetchChatChannelGroups(userID, w, r)
		case "createChatChannel":
			resultStruct = createChatChannel(userID, w, r)
		case "fetchChatChannels":
			resultStruct = fetchChatChannels(userID, w, r)
		case "removeChatChannelGroup":
			resultStruct = removeChatChannelGroup(userID, w, r)
		case "removeChatChannel":
			resultStruct = removeChatChannel(userID, w, r)
		case "addUsersToChannel":
			resultStruct = addUsersToChannel(userID, sessionId, w, r)
		case "fetchChatChannelUsers":
			resultStruct = fetchChatChannelUsers(userID, w, r)
		case "removeChatChannelUser":
			resultStruct = removeChatChannelUser(userID, sessionId, w, r)
		default:
		}
		marshaled, _ := json.Marshal(resultStruct)
		result = string(marshaled)
	}
	return
}

func sendMessage(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result structs.ResultMessage) {
	var parameters struct {
		From         *string
		To           *string
		Type         *string
		Id           *string
		Text         *string
		Participants *[]int64
		Attachments  *[]structs.ChatAttachment
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, message, err := authorization.SendMessage(userID, sessionId, parameters.From, parameters.To, parameters.Participants, parameters.Type, parameters.Id, parameters.Text, parameters.Attachments)
	if success && err == nil {
		result.Success = true
		result.Message = message
	} else if err != nil {
		result.Error = err.Error()
	}
	return result
}

func sendReplyMessage(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result structs.ResultMessage) {
	var parameters struct {
		BatchId                  *string
		From                     *string
		To                       *string
		RepliesToMessageId       *string
		RepliesToMessageBatchId  *string
		RepliesToMessageBatchId1 *string
		RepliesToMessageType     *string
		Id                       *string
		Text                     *string
		Attachments              *[]structs.ChatAttachment
		Participants             *[]string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, message, err := authorization.SendReplyMessage(userID, sessionId,
		parameters.BatchId, parameters.From, parameters.To, parameters.RepliesToMessageId,
		parameters.RepliesToMessageBatchId, parameters.RepliesToMessageBatchId1, parameters.RepliesToMessageType,
		parameters.Id, parameters.Text, parameters.Attachments, parameters.Participants)
	if success && err == nil {
		result.Success = true
		result.Message = message
	} else if err != nil {
		result.Error = err.Error()
	}
	return result
}

func receiveMessage(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ResultValue) {
	var parameters struct {
		From        *string
		To          *string
		Type        *string
		Id          *string
		Text        *string
		Time        int64
		BatchId     *string
		Attachments *[]structs.ChatAttachment
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, batchId, err := authorization.ReceiveMessage(userID, parameters.From, parameters.To, parameters.Type, parameters.Id, parameters.Text, parameters.Time, parameters.BatchId, parameters.Attachments)
	if success && err == nil {
		result.Success = true
		result.Value = batchId
	} else if err != nil {
		result.Error = err.Error()
	}
	return result
}

func getChatState(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result structs.ChatStateResult) {
	state, errS := authorization.GetChatState(userID, sessionId)
	if errS == nil {
		result.ChatState = *state
		result.Success = true
	} else {
		result.Error = errS.Error()
	}
	return result
}

func fetchChat(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ChatMessageResult) {
	var parameters struct {
		From          *string
		To            *string
		Type          *string
		Id            *string
		IncludeLatest bool
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	messages, participants, errM := authorization.FetchChat(userID, parameters.From, parameters.To, parameters.Type, parameters.Id, parameters.IncludeLatest)
	if errM == nil {
		result.Data = messages
		result.Participants = participants
		result.Success = true
	} else {
		result.Error = errM.Error()
	}
	return result
}

func fetchChatReplies(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ChatMessageResult) {
	var parameters struct {
		BatchId     *string
		MessageType *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)

	messages, partparameters, errM := authorization.FetchChatReplies(userID, parameters.BatchId, parameters.MessageType)
	if errM == nil {
		result.Data = messages
		result.Participants = partparameters
		result.Success = true
	} else {
		result.Error = errM.Error()
	}
	return result
}

func updateChat(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		BatchId        *string
		BatchId1       *string
		MessageId      *string
		MessageType    *string
		UpdateType     *string
		OpperationType *string
		Value          *string
		From           *string
		To             *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, errR := authorization.UpdateChat(userID, sessionId, parameters.BatchId, parameters.BatchId1, parameters.MessageId, parameters.MessageType, parameters.UpdateType, parameters.OpperationType, parameters.Value, parameters.From, parameters.To)
	if success && errR == nil {
		result.Success = true
	} else if errR != nil {
		result.Error = errR.Error()
	}
	return result
}

func updateReceiveChat(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		BatchId        *string
		MessageId      *string
		MessageType    *string
		UpdateType     *string
		OpperationType *string
		Value          *string
		From           *string
		To             *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, errR := authorization.UpdateReceiveChat(userID, sessionId, parameters.BatchId, parameters.MessageId, parameters.MessageType, parameters.UpdateType, parameters.OpperationType, parameters.Value, parameters.From, parameters.To)
	if success && errR == nil {
		result.Success = true
	} else if errR != nil {
		result.Error = errR.Error()
	}
	return result
}

func updateChatReply(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		BatchId        *string
		MessageId      *string
		UpdateType     *string
		MessageType    *string
		OpperationType *string
		Value          *string
		From           *string
		To             *string
		Participants   *[]string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, errR := authorization.UpdateChatReply(userID, sessionId, parameters.BatchId, parameters.MessageId, parameters.MessageType, parameters.UpdateType, parameters.OpperationType, parameters.Value, parameters.From, parameters.To, parameters.Participants)
	if success && errR == nil {
		result.Success = true
	} else if errR != nil {
		result.Error = errR.Error()
	}
	return result
}

func downloadAttachment(userID *int64, w *http.ResponseWriter, r *http.Request) {
	Id := r.URL.Query().Get("id")
	FileName := r.URL.Query().Get("fileName")

	file, err := authorization.GetFileBody(userID, &Id)

	if err == nil && file != nil {
		(*w).Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", FileName))
		(*w).Header().Set("Content-Type", "application/octet-stream")
		io.Copy((*w), *file)
	}
}

func markChatUnread(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		To   *string
		Type *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, errMU := authorization.MarkChatUnread(userID, parameters.To, parameters.Type)
	if success && errMU == nil {
		result.Success = true
	} else if errMU != nil {
		result.Error = errMU.Error()
	}
	return result
}

func initThread(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Label                *string
		ThreadId             *string
		Participants         *[]int64
		InstanceParticipants *[]int64
		ChatHistory          *structs.ChatHistory
		Message              *structs.Message
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, errMU := authorization.InitChatThread(userID, parameters.Label, parameters.ThreadId, parameters.Participants, parameters.InstanceParticipants, parameters.ChatHistory, parameters.Message)
	if success && errMU == nil {
		result.Success = true
	} else if errMU != nil {
		result.Error = errMU.Error()
	}
	return result
}

func receiveThread(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		ThreadId     *string
		Participants *[]int64
		ChatHistory  *structs.ChatHistory
		Message      *structs.Message
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, errMU := authorization.ReceiveThread(userID, parameters.ThreadId, parameters.Participants, parameters.ChatHistory, parameters.Message)
	if success && errMU == nil {
		result.Success = true
	} else if errMU != nil {
		result.Error = errMU.Error()
	}
	return result
}

func getThreads(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ResultValue) {
	var parameters struct {
		ThreadIds *[]string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	chatHistories, errCHH := authorization.GetThreadsHistories(userID, parameters.ThreadIds)
	valueBytes, _ := json.Marshal(chatHistories)
	if errCHH == nil {
		result.Success = true
		value := string(valueBytes[:])
		result.Value = &value
	} else {
		result.Success = false
		result.Error = errCHH.Error()
	}
	return result
}

func updateThread(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Category      *string
		ThreadId      *string
		Participants  *[]int64
		MessageUpdate *structs.UserMessageUpdate
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, errMU := authorization.UpdateThread(userID, parameters.Category, parameters.ThreadId, parameters.Participants, parameters.MessageUpdate)
	if success && errMU == nil {
		result.Success = true
	} else if errMU != nil {
		result.Error = errMU.Error()
	}
	return result
}

func createChatGroup(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Name *string
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, errMU := authorization.CreateChatGroup(userID, parameters.Name)
	if success && errMU == nil {
		result.Success = true
	} else if errMU != nil {
		result.Error = errMU.Error()
	}
	return result
}

func fetchChatChannelGroups(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ChatChannelGroupResult) {
	chatChannelGroups, errMU := authorization.FetchChatChannelGroups(userID)
	if errMU == nil {
		result.Success = true
		result.Data = &chatChannelGroups
	} else {
		result.Error = errMU.Error()
	}
	return result
}

func createChatChannel(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		Name    *string
		AmperId *int64
		GroupId *int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, errMU := authorization.CreateChatChannel(userID, parameters.Name, parameters.AmperId, parameters.GroupId)
	if success && errMU == nil {
		result.Success = true
	} else if errMU != nil {
		result.Error = errMU.Error()
	}
	return result
}

func fetchChatChannels(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.ChatChannelResult) {
	var parameters struct {
		GroupId *int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	chatChannels, errMU := authorization.FetchChatChannels(userID, parameters.GroupId)
	if errMU == nil {
		result.Success = true
		result.Data = &chatChannels
	} else {
		result.Error = errMU.Error()
	}
	return result
}

func removeChatChannelGroup(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		GroupId *int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, errRCCG := authorization.RemoveChatChannelGroup(userID, parameters.GroupId)
	if success && errRCCG == nil {
		result.Success = true
	} else if errRCCG != nil {
		result.Error = errRCCG.Error()
	}
	return result
}

func removeChatChannel(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		ChannelId *int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, errRCCG := authorization.RemoveChatChannel(userID, parameters.ChannelId)
	if success && errRCCG == nil {
		result.Success = true
	} else if errRCCG != nil {
		result.Error = errRCCG.Error()
	}
	return result
}

func addUsersToChannel(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		ChannelId *int64
		UserIds   *[]int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, errMU := authorization.AddUsersToChannel(userID, sessionId, parameters.ChannelId, parameters.UserIds)
	if success && errMU == nil {
		result.Success = true
	} else if errMU != nil {
		result.Error = errMU.Error()
	}
	return result
}

func fetchChatChannelUsers(userID *int64, w *http.ResponseWriter, r *http.Request) (result structs.UserReslt) {
	var parameters struct {
		ChannelId *int64
		Search    *[]string
		Start     int
		Limit     int
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	users, resultTotalCount, errMU := authorization.FetchChatChannelUsers(userID, parameters.ChannelId, parameters.Search, parameters.Start, parameters.Limit)
	if errMU == nil {
		result.Data = *users
		result.TotalCount = int(resultTotalCount)
		result.Success = true
	} else if errMU != nil {
		result.Error = errMU.Error()
	}
	return result
}

func removeChatChannelUser(userID *int64, sessionId *string, w *http.ResponseWriter, r *http.Request) (result structs.Result) {
	var parameters struct {
		ChannelId *int64
		UserIds   *[]int64
	}
	json.NewDecoder(r.Body).Decode(&parameters)
	success, errRCCU := authorization.RemoveChatChannelUser(userID, sessionId, parameters.ChannelId, parameters.UserIds)
	if success && errRCCU == nil {
		result.Success = true
	} else if errRCCU != nil {
		result.Error = errRCCU.Error()
	}
	return result
}
