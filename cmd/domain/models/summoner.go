package models

type Summoner struct {
}

type SummonerStats struct {
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
