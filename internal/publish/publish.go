package publish

import (
	"database/sql"
	"time"
)

func ParsePubDate(pubDate string) (pubTime sql.NullTime) {
	if pubDate == "" {
		pubTime.Valid = false
		return
	}

	pubTime.Valid = true
	t, err := time.Parse("Mon, 11 Nov 2024 13:28:52 -0700", pubDate)
	if err != nil {
		pubTime.Time = time.Now()
	}

	pubTime.Time = t
	return pubTime
}
