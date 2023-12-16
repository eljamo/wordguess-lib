package domain

type GameRecord struct {
	ID       string `json:"id"`
	PlayerID string `json:"player_id"`
	Result   int    `json:"result"`
	Rounds   int    `json:"rounds"`
	Date     string `json:"date"`
}
