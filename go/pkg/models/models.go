package models

import "time"

// Texts text output format.
type Texts map[string]string

// TranslatedText localized text.
type TranslatedText struct {
	ID        int
	Key       string
	Language  string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Language supported language.
type Language struct {
	ID        string
	CreatedAt time.Time
}
