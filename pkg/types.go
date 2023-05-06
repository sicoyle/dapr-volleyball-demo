package types

type Game struct {
	ID         int    `json:"id"`
	Round      int    `json:"round"`
	Team1Name  string `json:"team1Name"`
	Team2Name  string `json:"team2Name"`
	Team1Score int    `json:"team1Score"`
	Team2Score int    `json:"team2Score"`
}

type GameRequest struct {
	ID int `json:"id"`
}
