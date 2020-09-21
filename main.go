package main

import (
	"fmt"
	"github.com/ilhamfs/shopeewebminar/component"
	"github.com/ilhamfs/shopeewebminar/repository"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Module struct {
	redis          component.Redis
	userSaveRepo   repository.UserSaveRepo
	gameConfigRepo repository.GameConfigRepo
}

func NewHandler(redis component.Redis, userSaveRepo repository.UserSaveRepo, gameConfigRepo repository.GameConfigRepo) *Module {
	return &Module{
		redis:          redis,
		userSaveRepo:   userSaveRepo,
		gameConfigRepo: gameConfigRepo,
	}
}

func (m *Module) GetMapConfig(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func (m *Module) GetCharacterSaveFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("user_id"))
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

	handler := NewHandler(redis, repository.NewUserSaveRepo(db), repository.NewGameConfigRepo(db))

	router := httprouter.New()
	router.GET("/save/:user_id", handler.GetCharacterSaveFile)
	router.GET("/map/config", handler.GetMapConfig)

	log.Fatal(http.ListenAndServe(":8080", router))
}
