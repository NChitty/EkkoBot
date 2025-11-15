package riot

import "encoding/json"

type ErrorResponse struct {
	Status *Status `json:"status"`
}

type Status struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type AccountByRiotIdResponse struct {
	PlayerUuid string `json:"puuid"`
	Name       string `json:"gameName"`
	Tagline    string `json:"tagLine"`
}

const (
	RANKED_SOLO_DUO = iota
	RANKED_FLEX
)

type QueueType int

func (q *QueueType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	switch str {
	case "RANKED_FLEX_SR":
		*q = RANKED_FLEX
	case "RANKED_SOLO_5x5":
		*q = RANKED_SOLO_DUO
	}
	return nil
}

func (q QueueType) String() string {
	switch q {
	case RANKED_SOLO_DUO:
		return "Ranked Solo/Duo"
	case RANKED_FLEX:
		return "Ranked Flex"
	default:
		return "Unknown"
	}
}

type QueueResponse struct {
	LeagueId     string    `json:"leagueId"`
	QueueType    QueueType `json:"queueType"`
	Tier         string    `json:"tier"`
	Rank         string    `json:"string"`
	PlayerUuid   string    `json:"puuid"`
	LeaguePoints int       `json:"leaguePoints"`
	Wins         int       `json:"wins"`
	Losses       int       `json:"losses"`
	Veteran      bool      `json:"veteran"`
	Inactive     bool      `json:"inactive"`
	FreshBlood   bool      `json:"freshBlood"`
	HotStreak    bool      `json:"hotStreak"`
}
