package repository

import (
	"database/sql"
	"errors"
	"github.com/ilhamfs/shopeewebminar/model"
)

type gameConfigRepo struct {
	db *sql.DB
}

type GameConfigRepo interface {
	Get(gameID string) (config *model.GameConfig, err error)
}

func NewGameConfigRepo(db *sql.DB) GameConfigRepo {
	return &gameConfigRepo{
		db: db,
	}
}

func (r *gameConfigRepo) getDestDB(gameConfig *model.GameConfig) []interface{} {
	return []interface{}{
		&gameConfig.ID,
		&gameConfig.GameID,
		&gameConfig.Config,
	}
}

func (r *gameConfigRepo) Get(gameID string) (config *model.GameConfig, err error) {
	configs := make([]*model.GameConfig, 0)
	rows, err := r.db.Query("SELECT id, game_id, config FROM game_config WHERE game_id = ?", gameID)

	if err != nil {
		return
	}

	for rows.Next() {
		userDB := &model.GameConfig{}
		err = rows.Scan(r.getDestDB(userDB)...)
		if err != nil {
			return
		}
		configs = append(configs, userDB)
	}

	if len(configs) == 0 {
		err = errors.New("config not exists")
		return
	}

	config = configs[0]

	return
}
