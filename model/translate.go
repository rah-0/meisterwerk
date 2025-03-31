package model

import (
	"time"
)

type Language struct {
	Uuid        string    // Language UUID
	FirstInsert time.Time // Timestamp of first insert
	LastUpdate  time.Time // Timestamp of last update
	Prefix      string    // e.g., "en-US"
	Lang        string    // e.g., "English"
	Title       string    // e.g., "English" (native name)
	Img         string    // e.g., "/static/img/flags/us.png"
	MonthsShort string    // e.g., "Jan,Feb,Mar,Apr,..."
}

type LanguageKey struct {
	Uuid        string // Key ID (UUID)
	FirstInsert time.Time
	LastUpdate  time.Time
	Value       string // Semantic key string (e.g., "hello")
}

type LanguageValue struct {
	Uuid            string // Value row ID
	FirstInsert     time.Time
	LastUpdate      time.Time
	UuidLanguage    string // FK to Language
	UuidLanguageKey string // FK to LanguageKey
	Value           string // Translated text
}
