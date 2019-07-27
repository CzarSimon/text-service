package models

import "time"

// Texts text output format.
type Texts map[string]string

// Language supported language.
type Language struct {
	ID        string
	CreatedAt time.Time
}

// TranslatedText localized text.
type TranslatedText struct {
	ID        int
	Key       string
	Language  string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TextGroup group of texts.
type TextGroup struct {
	ID        string
	CreatedAt time.Time
}
