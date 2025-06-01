package utils

import "strings"

func (app *models.App) ExtractKeywords(title, description string) []string {
	text := strings.ToLower(title + " " + description)
	words := strings.Fields(text)

	var keywords []string
	seen := make(map[string]bool)

	for _, word := range words {
		// Clean word and filter by length
		cleaned := strings.Trim(word, ".,!?;:\"'()[]{}*&^%$#@")
		if len(cleaned) >= 3 && !seen[cleaned] {
			keywords = append(keywords, cleaned)
			seen[cleaned] = true
		}
	}

	return keywords
}
