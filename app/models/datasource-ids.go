package models

import "time"

type DatasourceIDs struct {
	MangadexID      string    `json:"mangadex_id"        db:"mangadex_id"`
	MangadexUpdated time.Time `json:"mangadex_updated"   db:"mangadex_updated"`
	// AnilistID        string    `json:"anilist_id"         db:"anilist_id"`
	// AnilistUpdated   time.Time `json:"anilist_updated"    db:"anilist_updated"`
	MangaPlusID      string    `json:"manga_plus_id"      db:"manga_plus_id"`
	MangaPlusUpdated time.Time `json:"manga_plus_updated" db:"manga_plus_updated"`
}
