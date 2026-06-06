package business

import (
	"amper/common/structs"
	"amper/data/database"
	"amper/properties/application"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"

	csmap "github.com/gnazaryan/concurrent-swiss-map"
)

var amperIdLock = &sync.Mutex{}
var amperId *int64

func AmperId() *int64 {
	if amperId == nil {
		amperIdLock.Lock()
		defer amperIdLock.Unlock()
		if amperId == nil {
			config, errAP := application.Get()
			if errAP == nil {
				identifier := config.GetString("amper.id", "")
				if len(identifier) > 0 {
					identifierInt, errAI := strconv.ParseInt(identifier, 10, 64)
					if errAI == nil {
						amperId = &identifierInt
					} else {
						log.Print(fmt.Errorf("make sure the amper.id property is integer value in the application.properties file"))
					}
				} else {
					log.Print(fmt.Errorf("make sure the amper.id property is configured in the application.properties file"))
				}
			} else {
				log.Print(fmt.Errorf("make sure application.properties file exists in luncher directory"))
			}

		}
	}

	return amperId
}

var amperCache = csmap.Create[int64, *structs.Amper](
	// set the number of map shards. the default value is 32.
	csmap.WithShardCount[int64, *structs.Amper](100),

	// if don't set custom hasher, use the built-in maphash.
	csmap.WithCustomHasher[int64, *structs.Amper](func(key int64) uint64 {
		return uint64(key)
	}),

	// set the total capacity, every shard map has total capacity/shard count capacity. the default value is 0.
	csmap.WithSize[int64, *structs.Amper](10000))

func GetAmperInstance(instanceId int64) *structs.Amper {
	return amperCache.StoreComputeSinglton(instanceId, func(value *structs.Amper) *structs.Amper {
		if value == nil {
			instance, errI := database.GetInstance(&instanceId, true)
			if errI == nil {
				return instance
			}
		}
		return value
	})
}

func InvalidateAmperCache() {
	amperCache.Clear()
}

func SystemDirectory() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(path.Dir(ex), "amper", "system")
}
