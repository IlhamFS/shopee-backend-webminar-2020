package repository

import (
	"database/sql"
	"errors"
	"github.com/ilhamfs/shopeewebminar/model"
)

type userSaveRepo struct {
	db *sql.DB
}

type UserSaveRepo interface {
	Get(username string) (save *model.UserSave, err error)
}

func NewUserSaveRepo(db *sql.DB) UserSaveRepo {
	return &userSaveRepo{
		db: db,
	}
}

func (r *userSaveRepo) getDestDB(save *model.UserSave) []interface{} {
	return []interface{}{
		&save.ID,
		&save.Username,
		&save.Save,
	}
}

func (r *userSaveRepo) Get(username string) (config *model.UserSave, err error) {
	saves := make([]*model.UserSave, 0)
	rows, err := r.db.Query("SELECT id, user_id, save FROM user_save WHERE username = ?", username)

	if err != nil {
		return
	}

	for rows.Next() {
		userSaveDB := &model.UserSave{}
		err = rows.Scan(r.getDestDB(userSaveDB)...)
		if err != nil {
			return
		}
		saves = append(saves, userSaveDB)
	}

	if len(saves) == 0 {
		err = errors.New("save not exists")
		return
	}

	config = saves[0]

	return
}
