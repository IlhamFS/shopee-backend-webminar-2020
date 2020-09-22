package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/ilhamfs/shopeewebminar/component"
	"github.com/ilhamfs/shopeewebminar/repository"
	"github.com/ilhamfs/shopeewebminar/util"
	gocache "github.com/patrickmn/go-cache"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Module struct {
	localCache     *gocache.Cache
	redis          component.Redis
	userSaveRepo   repository.UserSaveRepo
	gameConfigRepo repository.GameConfigRepo
}

func NewHandler(localCache *gocache.Cache, redis component.Redis, userSaveRepo repository.UserSaveRepo, gameConfigRepo repository.GameConfigRepo) *Module {
	return &Module{
		localCache:     localCache,
		redis:          redis,
		userSaveRepo:   userSaveRepo,
		gameConfigRepo: gameConfigRepo,
	}
}

func (m *Module) GetMapConfig(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gameID := ps.ByName("game_id")
	key := fmt.Sprintf("map-config::game-id-%s", gameID)
	log.Println("calling local cache")
	configInterface, found := m.localCache.Get(key)
	var config string
	if found {
		config = configInterface.(string)
		util.WriteOKResponse(w, config)
		return
	}

	log.Println("calling redis")
	configInterface, err := m.redis.Get(key)
	if err == redis.ErrNil {
		config = configInterface.(string)
		m.localCache.Set(key, config, gocache.DefaultExpiration)
		util.WriteOKResponse(w, config)
		return
	}

	log.Println("calling database")
	gameConfigModel, err := m.gameConfigRepo.Get(gameID)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
	config = gameConfigModel.Config
	m.redis.Set(key, config)
	m.localCache.Set(key, config, gocache.DefaultExpiration)
	util.WriteOKResponse(w, config)
	return
}

func (m *Module) GetCharacterSaveFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("username"))
}

func main() {
	redis, err := component.InitializeRedis()

	if err != nil {
		panic(err)
	}

	db, err := component.InitializeDatabase()

	if err != nil {
		panic(err)
	}
	localCache := gocache.New(5*time.Minute, 10*time.Minute)

	handler := NewHandler(localCache, redis, repository.NewUserSaveRepo(db), repository.NewGameConfigRepo(db))

	router := httprouter.New()
	router.GET("/save/:username", handler.GetCharacterSaveFile)
	router.GET("/map/config/:game_id", handler.GetMapConfig)

	log.Fatal(http.ListenAndServe(":8080", router))
}
