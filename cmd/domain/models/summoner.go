package models

type Summoner struct {
	ID         int64
	Name       string
	TagLine    string
	PlayerUuid string
}

type SummonerStats struct {
	Summoner           Summoner
	FlexGamesPlayed    int
	FlexTier           string
	FlexRank           string
	FlexWins           int
	FlexLp             int
	SoloDuoGamesPlayed int
	SoloDuoTier        string
	SoloDuoRank        string
	SoloDuoWins        int
	SoloDuoLp          int
}
