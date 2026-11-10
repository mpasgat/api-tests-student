package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

// ContainsAdvertisementByTitle проверяет, что в списке есть объявление с нужным названием
func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	var titles []string
	for _, item := range items {
		titles = append(titles, item.Title)
	}

	require.Contains(t, titles, title, "Объявление не найдено")
}
