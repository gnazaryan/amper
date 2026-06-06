package structs

import (
	"sync"
	"time"
)

type ResultMessage struct {
	Result
	Message *Message `json:"message"`
}

type ChatAttachment struct {
	Id        *string `json:"id"`
	Name      *string `json:"name"`
	Directory *string `json:"directory"`
}

type Message struct {
	From         *string            `json:"from"`
	FromUser     *User              `json:"fromUser"`
	To           *string            `json:"to"`
	Id           *string            `json:"id"`
	Text         *string            `json:"text"`
	DateTime     int64              `json:"dateTime"`
	BatchId      *string            `json:"batchId"`
	BatchId1     *string            `json:"batchId1"`
	Reactions    map[string][]int64 `json:"reactions"`
	Replies      map[string]int     `json:"replies"`
	ReplyBatchId *string            `json:"replyBatchId"`
	Deleted      bool               `json:"deleted"`
	Attachments  *[]ChatAttachment  `json:"attachments"`
}

type ChatHistoryItem struct {
	Id   *string `json:"id"`
	Full bool    `json:"full"`
}

type ChatThreadHistory struct {
	Label          *string `json:"label"`
	ThreadId       *string `json:"threadId"`
	UnreadMessages int     `json:"unreadMessages"`
	LastUpdateTime int64   `json:"lastUpdateTime"`
}

type ChatHistory struct {
	From           *string           `json:"from"`
	To             *string           `json:"to"`
	HistoryItems   []ChatHistoryItem `json:"historyItems"`
	LastUpdateTime int64             `json:"lastUpdateTime"`
	UnreadMessages int               `json:"unreadMessages"`
	Participants   *[]int64          `json:"participants"`
}

type ChatDirectItem struct {
	User        *User        `json:"user"`
	ChatHistory *ChatHistory `json:"chatHistory"`
}

type ChatThreadItem struct {
	Label       *string      `json:"label"`
	Users       *[]User      `json:"users"`
	ChatHistory *ChatHistory `json:"chatHistory"`
}

type ChatChannelItem struct {
	ChannelId   int64        `json:"channelId"`
	AmperId     int64        `json:"amperId"`
	Label       *string      `json:"label"`
	ChatHistory *ChatHistory `json:"chatHistory"`
}

type ChatChannelsGroup struct {
	GroupId  int64              `json:"groupId"`
	Label    *string            `json:"label"`
	Channels *[]ChatChannelItem `json:"channels"`
}

type ChatState struct {
	Directs              *[]ChatDirectItem    `json:"directs"`
	Threads              *[]ChatThreadItem    `json:"threads"`
	ChannelChannelGroups *[]ChatChannelsGroup `json:"channelGroups"`
}

type ChatStateResult struct {
	Result
	ChatState
}

type ChatHistoryResult struct {
	Result
	ChatHistorys *[]ChatHistory `json:"chatHistorys"`
}

type ChatMessageResult struct {
	Result
	Data         []Message       `json:"data"`
	Participants map[int64]*User `json:"participants"`
}

type MessageUpdate struct {
	MessageType    *string `json:"messageType"`
	UpdateType     *string `json:"updateType"`
	OpperationType *string `json:"opperationType"`
	From           *string `json:"from"`
	To             *string `json:"to"`
	Value          *string `json:"value"`
}

type UserUpdatesResult struct {
	Result
	Data map[string][]*interface{} `json:"data"`
}

type UserMessageUpdate struct {
	Message        *Message     `json:"message"`
	MessageType    *string      `json:"messageType"`
	UpdateType     *string      `json:"updateType"`
	OpperationType *string      `json:"opperationType"`
	From           *string      `json:"from"`
	To             *string      `json:"to"`
	Value          *string      `json:"value"`
	ChatHistory    *ChatHistory `json:"chatHistory"`
	Users          *[]User      `json:"users"`
}

type ChatChannelGroup struct {
	Id   *string `json:"id"`
	Name *string `json:"name"`
}

type ChatChannelGroupResult struct {
	Result
	Data *[]ChatChannelGroup `json:"data"`
}

type ChatChannel struct {
	Id              *int64            `json:"id"`
	Name            *string           `json:"name"`
	AmperId         *int64            `json:"amperId"`
	AmperName       *string           `json:"amperName"`
	GroupId         *int64            `json:"groupId"`
	GroupName       *string           `json:"groupName"`
	UserIds         *string           `json:"userIds"`
	BatchIds        *string           `json:"batchIds"`
	BatchIdsArray   []ChatHistoryItem `json:"batchIdsArray"`
	UserIdsInt64    *[]int64          `json:"userIdsInt64"`
	InstanceUserIDs map[int64][]int64 `json:"userIdsInt64"`
}

type ChatChannelResult struct {
	Result
	Data *[]ChatChannel `json:"data"`
}

type TimedMessages struct {
	Messages *[]Message
	Time     time.Time
}

type TimedMutex struct {
	Mutex *sync.RWMutex
	Time  time.Time
}

type TimedBatchId struct {
	BatchId *string
	Time    time.Time
}
