package utils

import "strings"

func (app *models.App) ArticleMatchesAlert(article Article, alert UserAlert) bool {
	// Check if any of the article keywords match any of the alert keywords
	for _, alertKeyword := range alert.Keywords {
		alertKeyword = strings.ToLower(strings.TrimSpace(alertKeyword))
		for _, articleKeyword := range article.Keywords {
			if strings.Contains(articleKeyword, alertKeyword) {
				return true
			}
		}
	}
	return false
}
