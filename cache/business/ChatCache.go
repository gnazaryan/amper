package business

import (
	"amper/common/structs"
	"amper/common/util"
	"amper/data/database"
	"encoding/json"
	"strconv"
	"strings"

	csmap "github.com/gnazaryan/concurrent-swiss-map"
)

var channelCache = csmap.Create[int64, *structs.ChatChannel](
	// set the number of map shards. the default value is 32.
	csmap.WithShardCount[int64, *structs.ChatChannel](100),

	// if don't set custom hasher, use the built-in maphash.
	csmap.WithCustomHasher[int64, *structs.ChatChannel](func(key int64) uint64 {
		return uint64(key)
	}),

	// set the total capacity, every shard map has total capacity/shard count capacity. the default value is 0.
	csmap.WithSize[int64, *structs.ChatChannel](10000))

func GetChatChannel(channelId *int64) *structs.ChatChannel {
	return channelCache.StoreComputeSinglton(*channelId, func(value *structs.ChatChannel) *structs.ChatChannel {
		if value == nil {
			chatChannel, errCC := database.FetchChatChannel(channelId)
			if errCC != nil {
				util.Loggify(errCC)
				return nil
			}
			chatChannel.BatchIdsArray = []structs.ChatHistoryItem{}
			json.Unmarshal([]byte(*chatChannel.BatchIds), &chatChannel.BatchIdsArray)

			channelUsersSplit := strings.Split(*chatChannel.UserIds, "__")
			var userIds []int64 = make([]int64, 0)
			var instanceUserIds map[int64][]int64 = map[int64][]int64{}
			for _, channelUserDirty := range channelUsersSplit {
				userId, _ := strconv.ParseInt(strings.Replace(channelUserDirty, "_", "", -1), 10, 64)
				userIds = append(userIds, userId)

				user := GetUser(&userId, true)
				if user != nil {
					if instanceUserIds[*user.AmperId] == nil {
						instanceUserIds[*user.AmperId] = []int64{userId}
					} else {
						instanceUserIds[*user.AmperId] = append(instanceUserIds[*user.AmperId], userId)
					}
				}
			}
			chatChannel.UserIdsInt64 = &userIds
			chatChannel.InstanceUserIDs = instanceUserIds
			value = chatChannel
		}
		return value
	})
}

func InvalidateChatChannel(channelId *int64) {
	channelCache.Delete(*channelId)
}
