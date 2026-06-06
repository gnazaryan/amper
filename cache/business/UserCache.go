package business

import (
	"amper/common/structs"
	"amper/common/util"
	"amper/data/database"

	csmap "github.com/gnazaryan/concurrent-swiss-map"
)

var userCache = csmap.Create[int64, *structs.User](
	// set the number of map shards. the default value is 32.
	csmap.WithShardCount[int64, *structs.User](1000),

	// if don't set custom hasher, use the built-in maphash.
	csmap.WithCustomHasher[int64, *structs.User](func(key int64) uint64 {
		return uint64(key)
	}),

	// set the total capacity, every shard map has total capacity/shard count capacity. the default value is 0.
	csmap.WithSize[int64, *structs.User](10000))

func GetUser(userId *int64, sensitive bool) *structs.User {
	return userCache.StoreComputeSinglton(*userId, func(value *structs.User) *structs.User {
		if value == nil {
			user, errDb := database.GetUser(userId, nil, nil, true, util.PointerBoolean(false), true, true)
			if errDb == nil {
				return user
			}
		}
		return value
	})
}

func InvalidateUserCache(userId *int64) {
	userCache.Delete(*userId)
}
