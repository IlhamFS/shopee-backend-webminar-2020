package main

import (
	"fmt"
	"github.com/ilhamfs/shopeewebminar/bloom"

	"github.com/ilhamfs/shopeewebminar/db"
	"github.com/ilhamfs/shopeewebminar/redis"
	"github.com/ilhamfs/shopeewebminar/repository"
	"github.com/ilhamfs/shopeewebminar/util"
	gocache "github.com/patrickmn/go-cache"
	"github.com/spaolacci/murmur3"
	"hash"
	"hash/fnv"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Module struct {
	bloomFilter    *bloom.BloomFilter
	localCache     *gocache.Cache
	redis          redis.RedisInterface
	userSaveRepo   repository.UserSaveRepo
	gameConfigRepo repository.GameConfigRepo
}

func NewHandler(bloomFilter *bloom.BloomFilter, localCache *gocache.Cache, redis redis.RedisInterface, userSaveRepo repository.UserSaveRepo,
	gameConfigRepo repository.GameConfigRepo) *Module {
	return &Module{
		bloomFilter:    bloomFilter,
		localCache:     localCache,
		redis:          redis,
		userSaveRepo:   userSaveRepo,
		gameConfigRepo: gameConfigRepo,
	}
}

// Local cache suitable in this API, because this API will be used intensively and using the same data across users
func (m *Module) GetMapConfig(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gameID := ps.ByName("game_id")

	log.Println("calling database")
	gameConfigModel, err := m.gameConfigRepo.Get(gameID)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
	}

	config := gameConfigModel.Config

	util.WriteOKResponse(w, config)
	return

}

// If we want to save cache space and load to avoid one time user, use bloom filter
func (m *Module) GetCharacterSaveFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")
	key := fmt.Sprintf("user-save::username-%s", username)

	log.Println("calling redis")
	var save string
	save, err := m.redis.Get(key)
	if err == nil {
		util.WriteOKResponse(w, save)
		return
	}

	log.Println("calling database")
	userSaveModel, err := m.userSaveRepo.Get(username)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
	save = userSaveModel.Save

	util.WriteOKResponse(w, save)
	return
}

// Create a web based battle royal
// We need to create two APIs:
// 1. Get map config for every user that join the game
// 2. Get user's save file
// Assume there's only one map in this game and used by a lot of users (High traffic)
func main() {
	// wrapper for redis
	redis, err := redis.InitializeRedis()

	if err != nil {
		panic(err)
	}
	// wrapper for database
	db, err := db.InitializeDatabase()

	if err != nil {
		panic(err)
	}
	// go-cache
	localCache := gocache.New(5*time.Minute, 10*time.Minute)
	// k = 3, we are using 3 hash functions
	// libraries available, redis also have an implementation of bloom filter
	bloomFilter := bloom.NewBloomFilter(100, []hash.Hash64{murmur3.New64(), fnv.New64(), fnv.New64a()})

	handler := NewHandler(bloomFilter, localCache, redis, repository.NewUserSaveRepo(db), repository.NewGameConfigRepo(db))
	router := httprouter.New()
	// ex: curl --location --request GET 'http://localhost:8080/character/save/ilhamfs'
	router.GET("/character/save/:username", handler.GetCharacterSaveFile)
	// ex: curl --location --request GET 'http://localhost:8080/map/config/the_dust'
	router.GET("/map/config/:game_id", handler.GetMapConfig)

	log.Fatal(http.ListenAndServe(":8080", router))
}
