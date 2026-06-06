package business

import (
	csmap "github.com/gnazaryan/concurrent-swiss-map"
)

var updates = csmap.Create[int64, map[string][]*interface{}](
	// set the number of map shards. the default value is 32.
	csmap.WithShardCount[int64, map[string][]*interface{}](64),

	// if don't set custom hasher, use the built-in maphash.
	csmap.WithCustomHasher[int64, map[string][]*interface{}](func(key int64) uint64 {
		return uint64(key)
	}),

	// set the total capacity, every shard map has total capacity/shard count capacity. the default value is 0.
	csmap.WithSize[int64, map[string][]*interface{}](1000))

func FetchUpdates(userID *int64) (result map[string][]*interface{}, err error) {
	var userUpdateCache = updates.DeleteRetrieve(*userID)
	return userUpdateCache, nil
}

func PutUpdate(userID *int64, category *string, value *interface{}) {

	updates.StoreCompute(*userID, func(valueMap map[string][]*interface{}) map[string][]*interface{} {
		if valueMap == nil {
			valueMap = make(map[string][]*interface{})
		}
		if valueMap[*category] == nil {
			valueMap[*category] = make([]*interface{}, 0)
		}
		valueMap[*category] = append(valueMap[*category], value)
		return valueMap
	})
}

func PutUpdates(userID *int64, Category *string, Participants *[]int64, Value *interface{}) (success bool, err error) {
	for _, participant := range *Participants {
		PutUpdate(&participant, Category, Value)
	}
	return true, nil
}
