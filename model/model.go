package model

type GameConfig struct {
	ID     string `json:"id"`
	GameID string `json:"game_id"`
	Config string `json:"config"`
}

type UserSave struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Save     string `json:"save"`
}
