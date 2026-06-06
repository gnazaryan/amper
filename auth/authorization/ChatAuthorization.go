package authorization

import (
	"amper/common/argument"
	"amper/common/structs"
	"amper/service/business"
	"io"
)

func SendMessage(userID *int64, sessionId *string, from *string, to *string, participants *[]int64, Type *string, id *string, text *string, Attachments *[]structs.ChatAttachment) (success bool, message *structs.Message, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "from": from, "to": to, "Type": Type, "id": id})
	if err != nil {
		return false, nil, err
	}
	return business.SendMessage(userID, sessionId, from, to, participants, Type, id, text, Attachments)
}

func SendReplyMessage(userID *int64, sessionId *string, BatchId *string, From *string, To *string, RepliesToMessageId *string, RepliesToMessageBatchId *string, RepliesToMessageBatchId1 *string, RepliesToMessageType *string, Id *string, Text *string, Attachments *[]structs.ChatAttachment, Participants *[]string) (success bool, message *structs.Message, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "from": From, "to": To, "repliesToMessageId": RepliesToMessageId, "repliesToMessageBatchId": RepliesToMessageBatchId, "repliesToMessageType": RepliesToMessageType, "id": Id})
	if err != nil {
		return false, nil, err
	}
	return business.SendReplyMessage(userID, sessionId, BatchId, From, To, RepliesToMessageId, RepliesToMessageBatchId, RepliesToMessageBatchId1, RepliesToMessageType, Id, Text, Attachments, Participants)
}

func ReceiveMessage(userID *int64, from *string, to *string, Type *string, id *string, text *string, time int64, batchId *string, Attachments *[]structs.ChatAttachment) (success bool, resultBatchId *string, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "from": from, "to": to, "Type": Type, "id": id, "time": time})
	if err != nil {
		return false, nil, err
	}
	return business.ReceiveMessage(userID, from, to, Type, id, text, time, batchId, Attachments)
}

func GetChatState(userID *int64, sessionId *string) (result *structs.ChatState, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "sessionId": sessionId})
	if err != nil {
		return nil, err
	}
	return business.GetChatState(userID, sessionId)
}

func FetchChat(userID *int64, from *string, to *string, Type *string, id *string, includeLatest bool) (result []structs.Message, participants map[int64]*structs.User, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "from": from, "to": to, "Type": Type})
	if err != nil {
		return nil, nil, err
	}
	return business.FetchChat(userID, from, to, Type, id, includeLatest)
}

func FetchChatReplies(userID *int64, BatchId *string, MessageType *string) (result []structs.Message, participants map[int64]*structs.User, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "batchId": BatchId, "MessageType": MessageType})
	if err != nil {
		return nil, nil, err
	}
	return business.FetchChatReplies(userID, BatchId, MessageType)
}

func UpdateChat(userID *int64, sessionId *string, BatchId *string, BatchId1 *string, MessageId *string, MessageType *string, UpdateType *string, OpperationType *string, Value *string, From *string, To *string) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "batchId": BatchId, "messageId": MessageId, "messageType": MessageType, "updateType": UpdateType, "opperationType": OpperationType, "value": Value, "from": From, "to": To})
	if err != nil {
		return false, err
	}
	return business.UpdateChat(userID, sessionId, BatchId, BatchId1, MessageId, MessageType, UpdateType, OpperationType, Value, From, To)
}

func UpdateReceiveChat(userID *int64, sessionId *string, BatchId *string, MessageId *string, MessageType *string, UpdateType *string, OpperationType *string, Value *string, From *string, To *string) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "sessionId": sessionId, "batchId": BatchId, "messageId": MessageId, "messageType": MessageType, "opperationType": OpperationType, "value": Value, "from": From, "to": To})
	if err != nil {
		return false, err
	}
	return business.UpdateReceiveChat(userID, sessionId, BatchId, MessageId, MessageType, UpdateType, OpperationType, Value, From, To)
}

func UpdateChatReply(userID *int64, sessionId *string, BatchId *string, MessageId *string, MessageType *string, UpdateType *string, OpperationType *string, Value *string, From *string, To *string, Participants *[]string) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "batchId": BatchId, "messageId": MessageId, "updateType": UpdateType, "opperationType": OpperationType, "value": Value, "from": From, "to": To, "participants": Participants})
	if err != nil {
		return false, err
	}
	return business.UpdateChatReply(userID, sessionId, BatchId, MessageId, MessageType, UpdateType, OpperationType, Value, From, To, Participants)
}

func GetFileBody(userID *int64, id *string) (result *io.ReadCloser, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "id": id})
	if err != nil {
		return nil, err
	}
	return business.GetFileBody(userID, id)
}

func MarkChatUnread(userID *int64, to *string, chatType *string) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "to": to, "type": chatType})
	if err != nil {
		return false, err
	}
	return business.MarkChatUnread(userID, to, chatType)
}

func InitChatThread(userID *int64, Label *string, ThreadId *string, Participants *[]int64, InstanceParticipants *[]int64, ChatHistory *structs.ChatHistory, Message *structs.Message) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "label": Label, "threadId": ThreadId, "instanceParticipants": InstanceParticipants})
	if err != nil {
		return false, err
	}
	return business.InitChatThread(userID, Label, ThreadId, Participants, InstanceParticipants, ChatHistory, Message)
}

func ReceiveThread(userID *int64, ThreadId *string, Participants *[]int64, ChatHistory *structs.ChatHistory, Message *structs.Message) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "threadId": ThreadId, "Participants": Participants})
	if err != nil {
		return false, err
	}
	return business.ReceiveThread(userID, ThreadId, Participants, ChatHistory, Message)
}

func GetThreadsHistories(userID *int64, ThreadIds *[]string) (result *[]structs.ChatHistory, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID})
	if err != nil {
		return nil, err
	}
	return business.GetThreadsHistories(userID, ThreadIds)
}

func UpdateThread(userID *int64, Category *string, ThreadId *string, Participants *[]int64, MessageUpdate *structs.UserMessageUpdate) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "ThreadId": ThreadId, "Category": Category})
	if err != nil {
		return false, err
	}
	return business.UpdateThread(userID, Category, ThreadId, Participants, MessageUpdate)
}

func CreateChatGroup(userID *int64, Name *string) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "name": Name})
	if err != nil {
		return false, err
	}
	return business.CreateChatGroup(userID, Name)
}

func FetchChatChannelGroups(userID *int64) (result []structs.ChatChannelGroup, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID})
	if err != nil {
		return nil, err
	}
	return business.FetchChatChannelGroups(userID)
}

func CreateChatChannel(userID *int64, Name *string, AmperId *int64, GroupId *int64) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "name": Name, "amperId": AmperId, "groupdId": GroupId})
	if err != nil {
		return false, err
	}
	return business.CreateChatChannel(userID, Name, AmperId, GroupId)
}

func FetchChatChannels(userID *int64, GroupId *int64) (result []structs.ChatChannel, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID})
	if err != nil {
		return nil, err
	}
	return business.FetchChatChannels(userID, GroupId)
}

func RemoveChatChannelGroup(userID *int64, groupId *int64) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "groupId": groupId})
	if err != nil {
		return false, err
	}
	return business.RemoveChatChannelGroup(userID, groupId)
}

func RemoveChatChannel(userID *int64, channelId *int64) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "channelId": channelId})
	if err != nil {
		return false, err
	}
	return business.RemoveChatChannel(userID, channelId)
}

func AddUsersToChannel(userID *int64, sessionId *string, ChannelId *int64, UserIds *[]int64) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "channelId": ChannelId, "UserIds": UserIds})
	if err != nil {
		return false, err
	}
	return business.AddUsersToChannel(userID, sessionId, ChannelId, UserIds)
}

func FetchChatChannelUsers(userID *int64, ChannelId *int64, Search *[]string, Start int, Limit int) (result *[]structs.User, resultTotalCount int, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "channelId": ChannelId})
	if err != nil {
		return nil, 0, err
	}
	return business.FetchChatChannelUsers(userID, ChannelId, Search, Start, Limit)
}

func RemoveChatChannelUser(userID *int64, sessionId *string, ChannelId *int64, UserIds *[]int64) (success bool, err error) {
	err = argument.Validate(map[string]interface{}{"userID": userID, "channelId": ChannelId, "userIds": UserIds})
	if err != nil {
		return false, err
	}
	return business.RemoveChatChannelUser(userID, sessionId, ChannelId, UserIds)
}
