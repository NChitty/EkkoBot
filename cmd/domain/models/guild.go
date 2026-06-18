package models

import "time"

type Guild struct {
	ID          int64
	DiscordID   string
	LastUpdated time.Time
}
