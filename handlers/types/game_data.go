package types

import "time"

type WordData struct {
	WordId  uint     `json:"-"`
	Letters []string `json:"letters"`
}

type GameData struct {
	Length     int16     `json:"length"`
	MaxTrials  int8      `json:"max_trials"`
	Trials     int8      `json:"trials"`
	Guessed    int16     `json:"guessed"`
	Expiration time.Time `json:"expiration"`
}

type GameResponse struct {
	WordData
	GameData
}

type GuessMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type GuessResponse struct {
	IsValid bool         `json:"valid"`
	Trials  int8         `json:"trials"`
	Endgame bool         `json:"endgame"`
	Message GuessMessage `json:"msg"`
}
