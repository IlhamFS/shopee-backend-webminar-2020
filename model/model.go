package model

type GameConfig struct {
	ID     string `json:"id"`
	Config string `json:"config"`
}

type UserSave struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Save   string `json:"save"`
}
