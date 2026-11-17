package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites/models"
)

// ContainsFavorite проверяет, что объявление есть в списке избранного
func ContainsFavorite(t *testing.T, list models.FavoritesList, advertisementId string) {
	require.Containsf(t, favoriteAdvertisementIds(list), advertisementId,
		"Объявления %s нет в избранном", advertisementId)
}

// NotContainsFavorite проверяет, что объявления нет в списке избранного
func NotContainsFavorite(t *testing.T, list models.FavoritesList, advertisementId string) {
	require.NotContainsf(t, favoriteAdvertisementIds(list), advertisementId,
		"Объявление %s всё ещё в избранном", advertisementId)
}

// favoriteAdvertisementIds собирает id объявлений из списка избранного
func favoriteAdvertisementIds(list models.FavoritesList) []string {
	var ids []string
	for _, item := range list.Items {
		ids = append(ids, item.Advertisement.Id)
	}

	return ids
}
